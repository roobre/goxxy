package mangler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type RegexMangler struct {
	headerRegexes map[*regexp.Regexp][]regexpReplace
	bodyRegexes   []regexpReplace
}

type regexpReplace struct {
	Regexp  *regexp.Regexp
	Replace string
}

func (rm *RegexMangler) AddHeaderRegex(matchHeader, search, replace string) *RegexMangler {
	headerRegex := regexp.MustCompile(matchHeader)
	searchRegex := regexp.MustCompile(search)
	rm.headerRegexes[headerRegex] = append(rm.headerRegexes[headerRegex], regexpReplace{searchRegex, replace})

	return rm
}

func (rm *RegexMangler) AddBodyRegex(search, replace string) *RegexMangler {
	searchRegex := regexp.MustCompile(search)
	rm.bodyRegexes = append(rm.bodyRegexes, regexpReplace{searchRegex, replace})

	return rm
}

func (rm *RegexMangler) Mangle(response *http.Response) *http.Response {
	for nameRegex, valueRegexes := range rm.headerRegexes {
		for name, headers := range response.Header {
			// TODO: Paralelize this
			if nameRegex.MatchString(name) {
				for _, header := range headers {
					for _, valueRegex := range valueRegexes {
						// Assume no need to make a new string
						header = valueRegex.Regexp.ReplaceAllString(header, valueRegex.Replace)
					}
				}
			}
		}
	}

	// Check len since we're copying body here
	if len(rm.bodyRegexes) > 0 {
		fullBody, _ := ioutil.ReadAll(response.Body)
		for _, regex := range rm.bodyRegexes {
			log.Printf("Searching for %s", regex.Regexp.String())
			fullBody = regex.Regexp.ReplaceAll(fullBody, []byte(regex.Replace))
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(fullBody))
	}

	return response
}
