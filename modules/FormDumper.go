package modules // import "roob.re/goxxy/modules"

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FormDumper logs to its Output request and response fields if they match their rules
type FormDumper struct {
	keywordSets        []keywordSet
	TryhardJson        bool
	IgnoreResponseCode bool
	Output             io.Writer
	maxSizer
}

type keywordSet struct {
	Type     uint8
	Keywords map[string]struct{}
}

const (
	keysetAny = iota
	keysetAll
)

// Add a set of keywords which will be compared against the requests and responses passing through the mangler. Those which contain all the keywords will be logged
func (d *FormDumper) All(keywords ...string) {
	d.add(keysetAll, keywords)
}

// Add a set of keywords which will be compared against the requests and responses passing through the mangler. Those which contain any of the keywords will be logged
func (d *FormDumper) Any(keywords ...string) {
	d.add(keysetAny, keywords)
}

func (d *FormDumper) add(t uint8, keywords []string) {
	var keywordsMap map[string]struct{}
	if len(keywords) > 0 {
		keywordsMap = make(map[string]struct{})
		for _, kw := range keywords {
			keywordsMap[kw] = struct{}{}
		}
	}

	d.keywordSets = append(d.keywordSets, keywordSet{Type: t, Keywords: keywordsMap})
}

func (d *FormDumper) Mangle(response *http.Response) *http.Response {
	if d.Output != nil && ((response.StatusCode < 400) || d.IgnoreResponseCode) {
		response.Request.ParseForm() // Idempotent

		// TODO Use reflect for this?
		keys := make(map[string]interface{})
		for key := range response.Request.Form {
			keys[key] = struct{}{}
		}

		if d.TryhardJson || strings.Contains(response.Header.Get("content-type"), "json") {
			if response.ContentLength <= d.maxSize() {
				buffer := CopyBody(response)
				json.Unmarshal(buffer, keys)
			}
		}

		var dump bool
		for i := 0; !dump && i < len(d.keywordSets); i++ {
			dump = shouldDump(&d.keywordSets[i], keys)
		}

		if dump {
			// TODO: Improve formatting
			d.Output.Write([]byte(fmt.Sprintln(response.Request.Form)))
		}
	}

	return response
}

func shouldDump(ks *keywordSet, keywords map[string]interface{}) bool {
	if len(keywords) < 1 {
		return false
	}

	if ks.Type == keysetAll {
		for wanted := range ks.Keywords {
			_, ok := keywords[wanted]
			if !ok {
				return false
			}
		}

		return true
	} else {
		for wanted := range ks.Keywords {
			_, ok := keywords[wanted]
			if ok {
				return true
			}
		}

		return false
	}
}
