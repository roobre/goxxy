package goxxy

import (
	"os"
	"roob.re/goxxy/modules"
	"testing"
)

func TestGoxxy(t *testing.T) {
	New()
	//	TODO
}

func TestDemux(t *testing.T) {
	// Create a new proxy
	proxy := New()

	// Add an EchoMangler, which just prints the request being mangled by this proxy
	proxy.AddMangler(modules.EchoMangler("Parent", os.Stdout))

	// Replace all links with https://www.roobre.es
	rm := &modules.RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, "https://www.roobre.es/")
	proxy.AddMangler(rm)

	// Add a new child, with their own matchers and manglers.
	// Requests are matched in depth, deepest match wins.
	// A proxy without any Matchers matches anything, but will prioritize their children if they match
	child1 := proxy.Child()
	child1.Match(HostMatcher(`(\w+\.)*google(\.\w{2,3})+`))
	child1.AddMangler(modules.EchoMangler("google anything:", os.Stdout))

	child11 := child1.Child()
	// Multiple Matchers are OR'ed together, if any of them matches, the proxy will mangle this request.
	child11.Match(HostMatcher(`google.es`))
	child11.Match(HostMatcher(`google.co.uk`))
	child11.AddMangler(modules.EchoMangler("google.es:", os.Stdout))

	child2 := proxy.Child()
	child2.Match(HostMatcher(`(\w+\.)*facebook\.\w{2,3}`))
	child2.AddMangler(modules.EchoMangler("facebook anything:", os.Stdout))
}
