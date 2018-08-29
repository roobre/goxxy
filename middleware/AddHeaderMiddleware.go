package middleware

import "net/http"

type HeaderAdder map[string]string

func (ha HeaderAdder) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for key, value := range ha {
			r.Header.Add(key, value)
		}

		handler.ServeHTTP(rw, r)
	})
}
