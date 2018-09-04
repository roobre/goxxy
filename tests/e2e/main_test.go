package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"roob.re/goxxy"
	"testing"
)

var g *goxxy.Goxxy
var proxy, upstream *httptest.Server

func TestMain(m *testing.M) {
	fmt.Println("Before tests")
	upstream = httptest.NewServer(http.HandlerFunc(testHandler))

	g = goxxy.New()
	proxy = httptest.NewServer(g)

	// Run tests
	r := m.Run()

	proxy.Close()
	upstream.Close()
	os.Exit(r)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.String()))
}

func TestA(t *testing.T) {
	if g == nil {
		t.Fail()
	}
}

func TestB(t *testing.T) {
	fmt.Println("B")
	t.Error("I always fail")
}
