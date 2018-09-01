package modules

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// HTMLModifier is anything capable of operating with a goquery.Document. Changes applied to the document will be reflected in the response set to the client
type HTMLModifier interface {
	ModifyHTML(doc *goquery.Document)
}

type HTMLModifierFunc func(doc *goquery.Document)

func (f HTMLModifierFunc) ModifyHTML(doc *goquery.Document) {
	f(doc)
}

// HTMLMangler is a convenience wrapper around goquery, which allows to apply Modifiers to the responses.
type HTMLMangler struct {
	modifiers []HTMLModifier
	maxSizer
}

func (h *HTMLMangler) AddModifier(modifier HTMLModifier) {
	h.modifiers = append(h.modifiers, modifier)
}

func (h *HTMLMangler) AddModifierFunc(modifier HTMLModifierFunc) {
	h.modifiers = append(h.modifiers, modifier)
}

func (h *HTMLMangler) Mangle(response *http.Response) *http.Response {
	if response.ContentLength > h.maxSize() {
		return response
	}

	if len(h.modifiers) <= 0 {
		return response
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Printf("error while building goquery document, response sent unmodified: %s\n", err.Error())
	}

	for _, modifier := range h.modifiers {
		modifier.ModifyHTML(document)
	}

	response.Body.Close()
	response.Body = ioutil.NopCloser(ioutil.NopCloser(strings.NewReader(document.Text())))
	return response
}
