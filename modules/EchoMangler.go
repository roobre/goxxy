package modules

import (
	"fmt"
	"io"
	"net/http"
	"roob.re/goxxy"
)

func EchoMangler(prefix string, w io.Writer) goxxy.Mangler {
	return goxxy.ManglerFunc(func(response *http.Response) *http.Response {
		w.Write([]byte(fmt.Sprintf("%s %s\n", prefix, response.Request.URL.String())))
		return response
	})
}
