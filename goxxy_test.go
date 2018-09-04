package goxxy

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestGoxxy(t *testing.T) {
	New()
	//	TODO
}

func TestDemux(t *testing.T) {
	nilBuffer := &bytes.Buffer{}

	req, err := http.NewRequest(http.MethodGet, "http://google.es", nilBuffer)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Create a new proxy
	proxy := New()
	// Add an EchoMangler, which just prints the request being mangled by this proxy

	// Test empty
	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(proxy).Pointer() {
		t.Error("Empty proxy did not match")
	}

	// Add a new child, with their own matchers and manglers.
	// Requests are matched in depth, deepest match wins.
	// A proxy without any Matchers matches anything, but will prioritize their children if they match
	child1 := proxy.Child()

	// Now we have a child, so it should match instead
	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(child1).Pointer() {
		t.Error("Empty child did not match")
	}

	// Add a matcher to it
	child1.Match(HostMatcher(`google\..+`))

	child11 := child1.Child()
	// Multiple Matchers are OR'ed together, if any of them matches, the proxy will mangle this request.
	child11.Match(HostMatcher(`google.es`))
	child11.Match(HostMatcher(`google.co.uk`))

	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(child11).Pointer() {
		t.Error("Level two child did not match")
	}

	req, _ = http.NewRequest(http.MethodGet, "http://google.com", nilBuffer)
	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(child1).Pointer() {
		t.Error("Level one child did not match")
	}

	child12 := child1.Child()
	// Multiple Matchers are OR'ed together, if any of them matches, the proxy will mangle this request.
	child12.Match(HostMatcher(`google.nonexistent`))

	req, _ = http.NewRequest(http.MethodGet, "http://google.nonexistent", nilBuffer)
	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(child12).Pointer() {
		t.Error("Level two child did not match")
	}

	child2 := proxy.Child()
	child2.Match(HostMatcher(`(\w+\.)*facebook\.\w{2,3}`))

	req, _ = http.NewRequest(http.MethodGet, "http://facebook.com", nilBuffer)
	if reflect.ValueOf(proxy.demux(req)).Pointer() != reflect.ValueOf(child2).Pointer() {
		t.Error("No match for facebook")
	}

	proxy.Match(HostMatcher("something"))

	if p := proxy.demux(req); reflect.ValueOf(p).Pointer() == reflect.ValueOf(proxy).Pointer() ||
		reflect.ValueOf(p).Pointer() != 0 {
		t.Error("Match for something that shouldn't")
	}
}
