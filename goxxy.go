package goxxy

import (
	"io"
	"log"
	"net/http"
)

var nopGoxxy = Goxxy{Client: http.DefaultClient}

type Middleware interface {
	Middleware(handler http.Handler) http.Handler
}
type MiddlewareFunc func(handler http.Handler) http.Handler

func (mf MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mf(handler)
}

type Mangler interface {
	Mangle(response *http.Response) *http.Response
}
type ManglerFunc func(response *http.Response) *http.Response

func (mf ManglerFunc) Mangle(response *http.Response) *http.Response {
	return mf(response)
}

type Matcher interface {
	Match(r *http.Request) bool
}
type MatcherFunc func(r *http.Request) bool

func (mf MatcherFunc) Match(r *http.Request) bool {
	return mf(r)
}

type Goxxy struct {
	Client      *http.Client
	ErrHandler  http.Handler
	middlewares []Middleware
	manglers    []Mangler
	matchers    []Matcher
}

func New() *Goxxy {
	return &Goxxy{Client: http.DefaultClient}
}

func (g *Goxxy) AddMiddleware(mw Middleware) {
	g.middlewares = append(g.middlewares, mw)
}
func (g *Goxxy) AddMiddlewareFunc(mw MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw)
}

func (g *Goxxy) AddMangler(mg Mangler) {
	g.manglers = append(g.manglers, mg)
}
func (g *Goxxy) AddManglerFunc(mg ManglerFunc) {
	g.manglers = append(g.manglers, mg)
}

func (g *Goxxy) Match(r *http.Request) http.Handler {
	if len(g.matchers) == 0 {
		return http.HandlerFunc(g.proxy)
	}

	for _, m := range g.matchers {
		if m.Match(r) {
			handler, isHandler := m.(http.Handler)
			if isHandler {
				return handler
			} else {
				return http.HandlerFunc(g.proxy)
			}
		}
	}

	return http.HandlerFunc(nopGoxxy.proxy)
}

func (g *Goxxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handler := g.Match(r)

	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i].Middleware(g)
	}

	handler.ServeHTTP(rw, r)
}

func (g *Goxxy) proxy(rw http.ResponseWriter, r *http.Request) {
	var url string
	if r.TLS != nil {
		url += "https://"
	} else {
		url += "http://"
	}

	url += r.Host + r.RequestURI

	log.Printf("Handling request to %s", url)
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
