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

	fd := &mangler.FormDumper{Output: os.Stdout}
	fd.All("user", "pwd")
	formdataFile, _ := os.Create("/tmp/formdata.txt")
	fd.Output = formdataFile
	goxxy.AddMangler(fd)

	http.ListenAndServe(":8080", goxxy)
}
