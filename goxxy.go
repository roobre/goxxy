package goxxy // import "roob.re/goxxy"

import (
	"io"
	"log"
	"net/http"
	"time"
)

var defaultClient = &http.Client{Timeout: 8 * time.Second, CheckRedirect: noRedirectsPolicy, Jar: nil}
var nopGoxxy = Goxxy{Client: defaultClient}

// Middleware is the de-facto standard interface for http middleware: Receives a handler, and returns another (typically a closure).
// Middlewares in Goxxy are used to modify a request before it is sent to the final server.
type Middleware interface {
	Middleware(handler http.Handler) http.Handler
}
type MiddlewareFunc func(handler http.Handler) http.Handler

func (mf MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mf(handler)
}

// Mangler is anything which can take an http.Response, do something with it, and then return it.
// Manglers which read Response.Body must care of leaving it untouched in the response they return, to ensure other
// manglers don't read partial responses.
type Mangler interface {
	Mangle(response *http.Response) *http.Response
}
type ManglerFunc func(response *http.Response) *http.Response

func (mf ManglerFunc) Mangle(response *http.Response) *http.Response {
	return mf(response)
}

// A module is anything which can operate both as a Middleware and as a Mangler
// This is, for now, unused. But I wanted to coin the term, for documentation readability purposes
type Module interface {
	Mangler
	Middleware
}

// Matcher is anything which can discern if a request should be intercepted or not
type Matcher interface {
	Match(*http.Request) bool
}
type MatcherFunc func(r *http.Request) bool

func (mf MatcherFunc) Match(r *http.Request) bool {
	return mf(r)
}

func noRedirectsPolicy(r *http.Request, rr []*http.Request) error {
	if len(rr) > 1 {
		return http.ErrUseLastResponse
	}
	return nil
}

// Goxxy is an http proxy which applies changes to requests and responses before and after sending them to the original server.
type Goxxy struct {
	Client          *http.Client // http.Client Goxxy will use to send requests upstream
	ErrHandler      http.Handler // ErrHandler will be invoked if the request made with Client fails with a non-recoverable error (e.g. NXDOMAIN, timeout, etc.)
	MangleRedirects bool
	middlewares     []Middleware
	manglers        []Mangler
	matchers        []Matcher
	children        []Goxxy
}

// New returns a fresh instance of Goxxy, with the default HTTP Client.
func New() *Goxxy {
	return &Goxxy{Client: defaultClient}
}

// AddMiddleware inserts a Module which will read and/or modify request before they are sent upstream
func (g *Goxxy) AddMiddleware(mw Middleware) {
	g.middlewares = append(g.middlewares, mw)
}

// AddMiddlewareFunc inserts a Module which will read and/or modify request before they are sent upstream
func (g *Goxxy) AddMiddlewareFunc(mw MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw)
}

// AddMangler inserts a Module which will read and/or modify responses after they're read from the target server and before they are sent back to the client
func (g *Goxxy) AddMangler(mg Mangler) {
	g.manglers = append(g.manglers, mg)
}

// AddManglerFunc inserts a Module which will read and/or modify responses after they're read from the target server and before they are sent back to the client
func (g *Goxxy) AddManglerFunc(mg ManglerFunc) {
	g.manglers = append(g.manglers, mg)
}

// Match adds a new matcher, which can discern if a request should be handled by this proxy or not. Multiple Matchers are OR'ed together.
// A Goxxy with no Matchers will match anything, but give priority to its children.
func (g *Goxxy) Match(m Matcher) {
	g.matchers = append(g.matchers, m)
}

// MatchFunc adds a new matcher, which can discern if a request should be handled by this proxy or not. Multiple Matchers are OR'ed together.
func (g *Goxxy) MatchFunc(m MatcherFunc) {
	g.matchers = append(g.matchers, m)
}

// Child creates adds a new child Goxxy and returns it.
func (g *Goxxy) Child() *Goxxy {
	g.children = append(g.children, Goxxy{Client: g.Client, ErrHandler: g.ErrHandler})
	return &g.children[len(g.children)-1]
}

// Goxxy won't follow redirects by default, since it can be breaking in some scenarios.
// However, enabling it will reduce latency and bandwith usage between goxxy and the clients.
func (g *Goxxy) FollowRedirects(follow bool) {
	if follow {
		g.Client.CheckRedirect = nil
	} else {
		g.Client.CheckRedirect = noRedirectsPolicy
	}
}

// Mangle returns a response after applying all manglers to the original one.
// Mangle is exported so Goxxy implements Mangler if needed, but it is not intended to be used from the outside in normal cases
func (g *Goxxy) Mangle(response *http.Response) *http.Response {
	// Do not invoke manglers if it's a redirect and MangleRedirects == false
	if g.MangleRedirects || !(response.StatusCode >= 300 && response.StatusCode < 400) {
		for i := range g.manglers {
			response = g.manglers[i].Mangle(response)
		}
	}

	return response
}

// Middleware returns the provided handler wrapped around g.middlewares
func (g *Goxxy) Middleware(handler http.Handler) http.Handler {
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i].Middleware(handler)
	}

	return handler
}

// ServeHTTP finds the appropiate Goxxy with demux(), wraps its proxy() with its Middleware() and calls it
func (g *Goxxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handlerGoxxy := g.demux(r)

	if handlerGoxxy == nil {
		log.Printf("Nothing matched `%s`, leaving intact", r.Method+" "+r.Host+r.RequestURI)
		nopGoxxy.proxy(rw, r)
		return
	}

	handlerGoxxy.Middleware(http.HandlerFunc(handlerGoxxy.proxy)).ServeHTTP(rw, r)
}

// Demux looks at matchers and children and returns a pointer to the Goxxy that should manage a given request
func (g *Goxxy) demux(r *http.Request) *Goxxy {
	var handler *Goxxy = nil

	if len(g.matchers) == 0 {
		if len(g.children) == 0 {
			// Return inmediately if empty
			return g
		}
		// Default as myself if no matchers
		handler = g
	} else {
		// Store myself if I match
		for _, c := range g.matchers {
			if c.Match(r) {
				handler = g
				break
			}
		}

		// Return if non-empty list of matchers and not matched any
		if handler == nil {
			return nil
		}
	}

	// Overwrite with children if they match
	for i := range g.children {
		child := &g.children[i]
		if childHandler := child.demux(r); childHandler != nil {
			handler = childHandler
			break
		}
	}

	return handler
}

// proxy makes a request to the upstream servers, mangles it, and echoes the response to the writer
func (g *Goxxy) proxy(rw http.ResponseWriter, r *http.Request) {
	var url string
	if r.TLS != nil {
		url += "https://"
	} else {
		url += "http://"
	}

	url += r.Host + r.RequestURI

	newreq, _ := http.NewRequest(r.Method, url, r.Body)
	newreq.Header = r.Header

	response, err := g.Client.Do(newreq)
	if err != nil {
		// Use custom handler if set
		if g.ErrHandler != nil {
			g.ErrHandler.ServeHTTP(rw, r)
			return
		}

		// This is a low-level error, so we just hijack the connection and forcefully close it
		if hijacker, isHijacker := rw.(http.Hijacker); isHijacker {
			conn, _, _ := hijacker.Hijack()
			conn.Close()
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("error during request: %v", err)
		return
	}

	copyResponse(rw, g.Mangle(response))
}

func copyResponse(rw http.ResponseWriter, response *http.Response) {
	for name, values := range response.Header {
		for _, value := range values {
			rw.Header().Add(name, value)
		}
	}

	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)
}
