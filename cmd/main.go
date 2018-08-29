package main

import (
	"net/http"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/mangler"
)

func main() {
	goxxy := goxxy.New()

	rm := &mangler.RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, "https://www.roobre.es/")
	goxxy.AddMangler(rm)

	formdataFile, _ := os.Create("/tmp/formdata.txt")
	fd := &mangler.FormDumper{Output: os.Stdout}
	fd.Output = formdataFile
	fd.All("user", "pwd")
	goxxy.AddMangler(fd)

	http.ListenAndServe(":8080", goxxy)
}
