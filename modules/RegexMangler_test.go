package modules

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"roob.re/goxxy/tests"
	"strings"
	"testing"
)

func TestRegexManglerHeadersRequest(t *testing.T) {
	rm := RegexMangler{}
	rm.AddHeaderRegex("User-Agent", "Firefox", "Roobreisafox")

	req := tests.Get()
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0")

	prevUa := req.Header.Get("User-Agent")

	rm.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != strings.Replace(prevUa, "Firefox", "Roobreisafox", 0) {
			t.Error("Replacing UA Failed")
		}
	}))
}

func TestRegexManglerHeadersResponse(t *testing.T) {
	const search = `^mock test (\w+)$`
	const replace = "Stable $1"

	rm := RegexMangler{}
	rm.AddHeaderRegex("Server", search, replace)

	resp := tests.GetResponse()
	origServer := resp.Header.Get("Server")

	resp = rm.Mangle(resp)

	if resp.Header.Get("Server") != regexp.MustCompile(search).ReplaceAllString(origServer, replace) {
		t.Error("Replace failed")
	}
}

func TestRegexManglerBodyResponse(t *testing.T) {
	const replacing = "https://www.roobre.es/"

	rm := RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, replacing)
	rm.AddBodyRegex(`([A-Z]\w*) webpage`, "$1 replacement")

	resp := rm.Mangle(tests.GetResponse())

	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)

	if strings.Count(string(buf.Bytes()), replacing) != 2 {
		t.Error("Link replace failed")
	}

	if strings.Count(string(buf.Bytes()), "Sample replacement") != 1 {
		t.Error("Text replace failed")
	}
}
