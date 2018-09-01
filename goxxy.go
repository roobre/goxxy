package goxxy

import (
	"io"
	"log"
	"net/http"
	"time"
)

var defaultClient = &http.Client{Timeout: 8 * time.Second}
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
//  manglers don't read partial responses.
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

// Goxxy is an http proxy which applies changes to requests and responses before and after sending them to the original server.
type Goxxy struct {
	Client      *http.Client
	ErrHandler  http.Handler
	middlewares []Middleware
	manglers    []Mangler
	matchers    []Matcher
	children    []Goxxy
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

func (g *Goxxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handler := g.demux(r)

	if handler == nil {
		log.Printf("Nothing matched `%s`, leaving intact", r.Method+" "+r.Host+r.RequestURI)
		nopGoxxy.proxy(rw, r)
		return
	}

	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i].Middleware(g)
	}

	handler.ServeHTTP(rw, r)
}

func (g *Goxxy) demux(r *http.Request) http.Handler {
	var handler http.Handler = nil

	if len(g.matchers) == 0 {
		if len(g.children) == 0 {
			// Return inmediately if empty
			return http.HandlerFunc(g.proxy)
		}
		// Default as myself if no matchers
		handler = http.HandlerFunc(g.proxy)
	}

	// Store myself if I match, noop for empty list
	for _, c := range g.matchers {
		if c.Match(r) {
			handler = http.HandlerFunc(g.proxy)
			break
		}
	}

	// Overwrite with children if they match
	for _, c := range g.children {
		if childHandler := c.demux(r); childHandler != nil {
			handler = childHandler
			break
		}
	}

	return handler
}

func (g *Goxxy) proxy(rw http.ResponseWriter, r *http.Request) {
	var url string
	if r.TLS != nil {
		url += "https://"
	} else {
		url += "http://"
	}

	url += r.Host + r.RequestURI

	//log.Printf("Handling request to %s", url)
	newreq, _ := http.NewRequest(r.Method, url, r.Body)

	response, err := g.Client.Do(newreq)
	if err != nil {
		// Use custom handler if set
		if g.ErrHandler != nil {
			g.ErrHandler.ServeHTTP(rw, r)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("error during request: %v", err)
		return
	}

	//for i := len(g.manglers) - 1; i >= 0; i-- {
	for i := range g.manglers {
		response = g.manglers[i].Mangle(response)
	}

	for name, values := range response.Header {
		for _, value := range values {
			rw.Header().Add(name, value)
		}
	}

	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)
}
