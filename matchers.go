package goxxy

import (
	"net/http"
	"regexp"
)

func HeaderMatcher(name, valueRegex string) Matcher {
	// TODO: Proper error log
	regex := regexp.MustCompile(valueRegex)
	return MatcherFunc(func(r *http.Request) bool {
		for name, value := range r.Header {
			if name == name && regex.MatchString(value[0]) {
				return true
			}
		}

		return false
	})
}

func HostMatcher(host string) Matcher {
	// TODO: Proper error log
	regex := regexp.MustCompile(host)
	return MatcherFunc(func(r *http.Request) bool {
		return regex.MatchString(r.Host)
	})
}
