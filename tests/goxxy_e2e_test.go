package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var upstream *httptest.Server
var client *http.Client

func TestMain(m *testing.M) {
	upstream = httptest.NewServer(http.HandlerFunc(HTMLHandler))
	client = upstream.Client()

	r := m.Run()

	upstream.Close()
	os.Exit(r)
}
