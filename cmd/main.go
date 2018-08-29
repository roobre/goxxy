package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/mangler"
)

func main() {
	goxxy := goxxy.New()

	rm := &mangler.RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, "https://www.roobre.es/")

	goxxy.AddMangler(rm)

	goxxy.AddManglerFunc(func(response *http.Response) *http.Response {
		body, _ := ioutil.ReadAll(response.Body)
		file, _ := os.Open("/tmp/whatever/" + response.Request.Host)
		file.Write(body)

		response.Body = ioutil.NopCloser(bytes.NewReader(body))
		return response
	})

	http.ListenAndServe(":8080", goxxy)
}
