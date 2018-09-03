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

func TestHTMLManglerEdgeCases(t *testing.T) {
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

	response := tests.GetResponse()
	prevResponse := response
	response = htmlMangler.Mangle(response)
	if response.Body != prevResponse.Body {
		t.Error("Response body unnecessarily modified, empty modifier list")
	}

	response.ContentLength = defaultResponseMaxSize + 2
	prevResponse = response
	response = htmlMangler.Mangle(response)
	htmlMangler.AddModifier(HTMLModifierFunc(linkReplacer))
	response = htmlMangler.Mangle(response)

	buf := bytes.Buffer{}
	io.Copy(&buf, response.Body)

	if response.Body != prevResponse.Body || bytes.Compare(buf.Bytes(), []byte(tests.ResponseHTML)) != 0 {
		t.Error("Response unmodified despite advertised content length")
	}
}
