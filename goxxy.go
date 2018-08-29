package goxxy

import (
	"io"
	"log"
	"net/http"
)

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

type Goxxy struct {
	Client      *http.Client
	middlewares []Middleware
	manglers    []Mangler
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

func (g *Goxxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(g.proxy)
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i].Middleware(g)
	}

	handler.ServeHTTP(rw, r)
}

func (g *Goxxy) proxy(rw http.ResponseWriter, r *http.Request) {
	var url string
	url = "http://"
	if r.Host != "" {
		url += r.Host
		//if parts := strings.Split(r.RemoteAddr, ":"); len(parts) > 1 {
		//	url += ":" + parts[1]
		//}
	} else {
		url += r.RemoteAddr
	}

	url += r.RequestURI

	log.Printf("Handling request to %s", url)
	newreq, _ := http.NewRequest(r.Method, url, r.Body)

	response, err := g.Client.Do(newreq)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("error during request: %v", err)
		return
	}

	for i := len(g.manglers) - 1; i >= 0; i-- {
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
