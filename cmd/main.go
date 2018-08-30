package main

import (
	"net/http"
	"os"
	"roob.re/goxxy"
	"roob.re/goxxy/modules"
)

func main() {
	goxxy := goxxy.New()

	rm := &modules.RegexMangler{}
	rm.AddBodyRegex(`https?://(?:\w+\.\w+)+/`, "https://www.roobre.es/")
	goxxy.AddMangler(rm)

	formdataFile, _ := os.Create("/tmp/formdata.txt")
	fd := &modules.FormDumper{Output: os.Stdout}
	fd.Output = formdataFile
	fd.All("user", "pwd")
	goxxy.AddMangler(fd)

	http.ListenAndServe(":8080", goxxy)
}
