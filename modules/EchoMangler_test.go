package modules

import (
	"bytes"
	"fmt"
	"roob.re/goxxy/tests"
	"testing"
)

func TestEchoMangler(t *testing.T) {
	buffer := bytes.Buffer{}
	mangler := EchoMangler("Testing", &buffer)

	response := tests.GetResponse()
	mangler.Mangle(response)

	if bytes.Compare(buffer.Bytes(), []byte(fmt.Sprintf("Testing %s\n", response.Request.URL.String()))) != 0 {
		t.Error("Log doesn't match")
	}
}
