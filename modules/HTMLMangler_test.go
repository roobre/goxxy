package modules

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"roob.re/goxxy/tests"
	"testing"
)

func TestHTMLMangler(t *testing.T) {
	linkReplacer := func(doc *goquery.Document) {
		for _, node := range doc.Find("a").Nodes {
			for i := range node.Attr {
				attr := &node.Attr[i]

				if attr.Key == "href" {
					attr.Val = "https://www.carrierlost.net/"
				}
			}
		}
	}

	htmlMangler := HTMLMangler{}
	htmlMangler.AddModifierFunc(linkReplacer)

	response := tests.GetResponse()
	response = htmlMangler.Mangle(response)

	buf := bytes.Buffer{}
	io.Copy(&buf, response.Body)

	if bytes.Contains(buf.Bytes(), []byte("example.org")) || !bytes.Contains(buf.Bytes(), []byte("https://www.carrierlost.net/")) {
		t.Error("Partial or not found replacement")
	}
}
