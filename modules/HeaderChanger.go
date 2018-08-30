package modules

import (
	"net/http"
	"strings"
)

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
