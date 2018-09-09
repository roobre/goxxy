package tests

import (
	"encoding/json"
	"io"
	"net/http"
)

var TestHandler = http.NewServeMux()

const RouteSimpleHTML = "/index.html"
const RouteSimpleJSON = "/index.json"
const ReouteRedirect = "/oldindex.html"
const RouteNotFound = "/notfound.html"
const RoutePostEcho = "/post.php"
const NetworkError = "/neterror"

func init() {
	TestHandler.HandleFunc(RouteSimpleHTML, func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(ResponseJSON)
	})

	TestHandler.HandleFunc(RouteSimpleJSON, func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(ResponseHTML))
	})

	TestHandler.HandleFunc(ReouteRedirect, func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, RouteSimpleHTML, http.StatusFound)
	})

	TestHandler.HandleFunc(RouteNotFound, func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "Not found", http.StatusNotFound)
	})

	TestHandler.HandleFunc(RoutePostEcho, func(rw http.ResponseWriter, r *http.Request) {
		io.Copy(rw, r.Body)
		rw.WriteHeader(http.StatusCreated)
	})

	TestHandler.HandleFunc(NetworkError, func(rw http.ResponseWriter, r *http.Request) {
		if hijacker, isHijacker := rw.(http.Hijacker); isHijacker {
			conn, _, _ := hijacker.Hijack()
			conn.Close()
			return
		} else {
			panic("ResponseWriter does not implement Hijacker")
		}
	})
}
