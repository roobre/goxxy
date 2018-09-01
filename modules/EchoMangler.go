package modules

import (
	"fmt"
	"io"
	"net/http"
	"roob.re/goxxy"
)

// EchoMangler returns a Mangler which echoes a string representation of the request being mangled to the supplied io.Writer
func EchoMangler(prefix string, w io.Writer) goxxy.Mangler {
	return goxxy.ManglerFunc(func(response *http.Response) *http.Response {
		w.Write([]byte(fmt.Sprintf("%s %s\n", prefix, response.Request.URL.String())))
		return response
	})
}
