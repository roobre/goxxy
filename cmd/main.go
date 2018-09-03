package main

import (
	"log"
	"net/http"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/modules"
)

func main() {
	// Create a new proxy
	proxy := goxxy.New()

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
	child1.Match(goxxy.HostMatcher(`(\w+\.)*google(\.\w{2,3})+`))
	child1.AddMangler(modules.EchoMangler("google anything:", os.Stdout))

	child11 := child1.Child()
	// Multiple Matchers are OR'ed together, if any of them matches, the proxy will mangle this request.
	child11.Match(goxxy.HostMatcher(`google.es`))
	child11.Match(goxxy.HostMatcher(`google.co.uk`))
	child11.AddMangler(modules.EchoMangler("google.es:", os.Stdout))

	child2 := proxy.Child()
	child2.Match(goxxy.HostMatcher(`(\w+\.)*facebook\.\w{2,3}`))
	child2.AddMangler(modules.EchoMangler("facebook anything:", os.Stdout))

	log.Println("Starting Goxxy on :8080")
	http.ListenAndServe(":8080", proxy)
}
