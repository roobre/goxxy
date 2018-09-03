package tests

import (
	"bytes"
	"fmt"
	"roob.re/goxxy/modules"
	"testing"
)

func TestEchoMangler(t *testing.T) {
	buffer := bytes.Buffer{}
	mangler := modules.EchoMangler("Testing", &buffer)

	response := GetResponse()
	mangler.Mangle(response)

	if bytes.Compare(buffer.Bytes(), []byte(fmt.Sprintf("Testing %s\n", response.Request.URL.String()))) != 0 {
		t.Error("Log doesn't match")
	}
}
