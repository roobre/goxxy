package modules

import (
	"bytes"
	"net/http"
	"roob.re/goxxy/tests"
	"testing"
)

func TestFormDumperShouldMatch(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.org/login?user=perry&password=platypus", &bytes.Buffer{})
	if err != nil {
		t.Error(err)
	}

	resp := tests.GetResponse()
	resp.Request = req

	out := &bytes.Buffer{}
	var prevLen int

	fd := FormDumper{Output: out}
	fd.Any("userino", "passworderino")

	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Match for non matching Any")
	}
	fd.keywordSets = nil

	fd.Any("passworderino", "user")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("Non match for a matching Any")
	}
	fd.keywordSets = nil

	fd.All("passworderino", "user")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Match for a single matching All")
	}
	fd.keywordSets = nil

	fd.All("password", "user")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("No match for all matching All")
	}
	fd.keywordSets = nil

	fd.All("u", "p")
	fd.All("user", "p")
	fd.All("password", "u")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Match for multiple non-satisfying All")
	}

	fd.Any("password")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("No match for complex ruleset including satisfying any")
	}
	fd.keywordSets = nil

	fd.All("u", "p")
	fd.Any("nope", "stillno")
	fd.All("user", "p")
	fd.All("password", "u")
	fd.Any("nono", "noisno")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Match for non-satisfying ruleset")
	}
	fd.keywordSets = nil

	fd.All("u", "p")
	fd.Any("nope", "stillno")
	fd.All("user", "p")
	fd.Any("nono", "noisno", "user")
	fd.All("password", "u")
	fd.Any("nono", "noisno")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("No match for satisfying ruleset")
	}
}

func TestFormDumperJson(t *testing.T) {
	resp := tests.GetResponseJSON()

	out := &bytes.Buffer{}
	var prevLen int

	fd := FormDumper{Output: out}
	fd.All("Name", "Count", "Value")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("No match for satisfying All")
	}

	resp.Header.Del("content-type")
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Shouldn't have looked at body due to missing Content-Type")
	}

	fd.TryhardJson = true
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("Should have looked at body due to TryhardJson = true")
	}
}

func TestFormDumperResponseCode(t *testing.T) {
	resp := tests.GetResponseJSON()
	resp.StatusCode = http.StatusForbidden

	out := &bytes.Buffer{}
	var prevLen int

	fd := FormDumper{Output: out}
	fd.All("Name", "Count", "Value")

	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() != prevLen {
		t.Error("Shouldn't have matched non-OK response")
	}

	fd.IgnoreResponseCode = true
	prevLen = out.Len()
	fd.Mangle(resp)
	if out.Len() == prevLen {
		t.Error("Should have matched non-OK response due to Ignore flag being true")
	}
}
