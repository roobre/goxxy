package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"roob.re/goxxy"
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

func TestIdentity(t *testing.T) {
	// TODO: Some TLS testing should be placed here, for now plain http testing is meaningless

	g := goxxy.New()
	g.Client = client

	proxy := httptest.NewServer(g)
	proxyClient := proxy.Client()

	req, _ := http.NewRequest(http.MethodGet, upstream.URL+"/example", &bytes.Buffer{})

	origResponse, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	proxyResponse, err := proxyClient.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err := CompareResponses(origResponse, proxyResponse); err != nil {
		t.Error(err)
	}
}
