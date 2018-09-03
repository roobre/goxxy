package tests

import (
	"bytes"
	"net/http"
	"roob.re/goxxy"
	"testing"
)

func TestHeaderMatcher(t *testing.T) {
	req := Get()
	req.Header.Set("X-Custom", "customvalue")

	var matched bool

	matched = goxxy.HeaderMatcher("User-Agent", ".*").Match(req)
	if !matched {
		t.Error("User-Agent header did not match")
	}

	matched = goxxy.HeaderMatcher("X-Custom", "^.*value$").Match(req)
	if !matched {
		t.Error("X-Custom header did not match")
	}

	matched = goxxy.HeaderMatcher("Nonexistant", ".*").Match(req)
	if matched {
		t.Error("Match for nonexistant header")
	}
	matched = goxxy.HeaderMatcher("Nonexistant", "").Match(req)
	if matched {
		t.Error("Match for nonexistant header")
	}
}

func TestHostMatcher(t *testing.T) {
	req := Get()

	var matched bool

	matched = goxxy.HostMatcher("example.org").Match(req)
	if !matched {
		t.Error("Did not match partial")
	}

	matched = goxxy.HostMatcher("^example.org$").Match(req)
	if matched {
		t.Error("Partial match on explicit regex")
	}

	req, _ = http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/asdf", &bytes.Buffer{})
	matched = goxxy.HostMatcher(`^127\.0\.0\.1$`).Match(req)
	if !matched {
		t.Error("Did not match with custom port not included in regex")
	}
}
