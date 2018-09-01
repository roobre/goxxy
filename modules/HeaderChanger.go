package modules

import (
	"net/http"
	"strings"
)

// HeaderChanger will add, append, or remove Headers before the request is sent to the server and before the response is sent to the client.
// Request is a map whose keys are the request headers to act on, and the values are the values to be added, removed, or appended.
// Response is a map whose keys are the request headers to act on, and the values are the values to be added, removed, or appended.
// For both maps, keys (header names) starting with "-" (e.g. "-Server") will cause the header to be deleted. Keys starting with "+", will cause the value to be appended to the header name, and keys without any prefix will set the value regardless of any previous value.
type HeaderChanger struct {
	Request  map[string]string
	Response map[string]string
}

func (ha *HeaderChanger) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ha.changeHeaders(ha.Request, r.Header)
		handler.ServeHTTP(rw, r)
	})
}

func (ha *HeaderChanger) Mangle(response *http.Response) *http.Response {
	ha.changeHeaders(ha.Response, response.Header)
	return response
}

func (ha *HeaderChanger) changeHeaders(changes map[string]string, headers http.Header) {
	for key, value := range ha.Request {
		if strings.HasPrefix(key, "-") {
			headers.Del(key)
		} else if strings.HasPrefix(key, "+") {
			headers.Add(key, value)
		} else {
			headers.Set(key, value)
		}
	}
}
