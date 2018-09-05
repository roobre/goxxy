package goxxy_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/tests"
	"testing"
)

var upstream *httptest.Server
var client *http.Client

func TestMain(m *testing.M) {
	upstream = httptest.NewServer(http.HandlerFunc(tests.HTMLHandler))
	client = upstream.Client()

	r := m.Run()

	upstream.Close()
	os.Exit(r)
}

func TestIdentity(t *testing.T) {
	// TODO: Some TLS testing should be placed here, for now plain http testing is meaningless

	g := goxxy.New()
	g.Client = client

	// Using a second httptest.Server doesnt count for coverage, apparently
	//proxy := httptest.NewServer(g)
	//proxyClient := proxy.Client()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, upstream.URL+"/example", &bytes.Buffer{})

	origResponse, err := client.Do(req)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	g.ServeHTTP(recorder, req)
	proxyResponse := recorder.Result()

	if err := tests.CompareResponses(origResponse, proxyResponse); err != nil {
		t.Error(err)
	}
}
