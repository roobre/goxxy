package modules // import "roob.re/goxxy/modules"

import (
	"net/http"
	"strings"
)

// HeaderChanger will add, append, or remove Headers before the request is sent to the server or before the response is sent to the client.
// Keys (header names) starting with "-" (e.g. "-Server") will cause the header to be deleted. Keys starting with "+", will cause the value to be appended to the header name, and keys without any prefix will set the value regardless of any previous value.
type HeaderChanger map[string]string

func (ha HeaderChanger) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ha.changeHeaders(r.Header)
		handler.ServeHTTP(rw, r)
	})
}

func (ha HeaderChanger) Mangle(response *http.Response) *http.Response {
	ha.changeHeaders(response.Header)
	return response
}

func (ha HeaderChanger) changeHeaders(headers http.Header) {
	for key, value := range ha {
		if strings.HasPrefix(key, "-") {
			headers.Del(key)
		} else if strings.HasPrefix(key, "+") {
			headers.Add(key, value)
		} else {
			headers.Set(key, value)
		}
	}
}
