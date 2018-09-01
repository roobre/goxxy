package modules

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// RegexMangler is a collection of regexes to apply to responses which will be set back to the client, both to the headers and body.
type RegexMangler struct {
	headerRegexes map[string][]regexpReplace
	bodyRegexes   []regexpReplace
	maxSizer
}

type regexpReplace struct {
	Regexp  *regexp.Regexp
	Replace string
}

// AddHeaderRegex adds a new regex which will be applied to the headers sent in the response. header is the header name and must match verbatim.
func (rm *RegexMangler) AddHeaderRegex(header, search, replace string) *RegexMangler {
	searchRegex := regexp.MustCompile(search)
	rm.headerRegexes[header] = append(rm.headerRegexes[header], regexpReplace{searchRegex, replace})

	return rm
}

// AddBodyRegex adds a new regex which will be applied to the response body.
func (rm *RegexMangler) AddBodyRegex(search, replace string) *RegexMangler {
	searchRegex := regexp.MustCompile(search)
	rm.bodyRegexes = append(rm.bodyRegexes, regexpReplace{searchRegex, replace})

	return rm
}

func (rm *RegexMangler) Mangle(response *http.Response) *http.Response {
	if response.ContentLength > rm.maxSize() {
		return response
	}

	for headerName, valueRegexes := range rm.headerRegexes {
		for name, headers := range response.Header {
			// TODO: Paralelize this
			if name == headerName {
				//for _, header := range headers { // Ignore multi-valued headers for now
				for _, valueRegex := range valueRegexes {
					headers[0] = valueRegex.Regexp.ReplaceAllString(headers[0], valueRegex.Replace)
				}
				//}
			}
		}
	}

	// Check len since we're copying body here
	if len(rm.bodyRegexes) > 0 {
		fullBody := CopyBody(response)

		for _, regex := range rm.bodyRegexes {
			log.Printf("Searching for %s", regex.Regexp.String())
			fullBody = regex.Regexp.ReplaceAll(fullBody, []byte(regex.Replace))
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(fullBody))
	}

	return response
}
