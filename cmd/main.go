package main

import (
	"log"
	"net/http"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/modules"
)

func main() {
	proxy := goxxy.New()

	proxy.AddMangler(modules.EchoMangler("Parent", os.Stdout))

	rm := &modules.RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, "https://www.roobre.es/")
	proxy.AddMangler(rm)

	child1 := proxy.Child()
	child1.Match(goxxy.HostMatcher(`(\w+\.)*google\.\w{2,3}`))
	child1.AddMangler(modules.EchoMangler("google anything:", os.Stdout))

	child11 := child1.Child()
	child11.Match(goxxy.HostMatcher(`google.es`))
	child11.AddMangler(modules.EchoMangler("google.es:", os.Stdout))

	child2 := proxy.Child()
	child2.Match(goxxy.HostMatcher(`(\w+\.)*facebook\.\w{2,3}`))
	child2.AddMangler(modules.EchoMangler("facebook anything:", os.Stdout))

	log.Println("Starting Goxxy on :8080")
	http.ListenAndServe(":8080", proxy)
}
