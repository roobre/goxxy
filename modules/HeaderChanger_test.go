package modules

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestHeaderChangerSet(t *testing.T) {
	headers := http.Header{}
	headers.Add("Already-Present", "1")

	changer := HeaderChanger{}
	changer["Added"] = "Asdf"

	changer.changeHeaders(headers)
	if headers.Get("Added") != "Asdf" {
		t.Error("Missing header wasn't added")
	}

	changer["Already-Present"] = "2"
	changer.changeHeaders(headers)
	if strings.Join(headers["Already-Present"], ", ") != "2" {
		t.Errorf("Existing header wasn't replaced: %v", headers["Already-Present"])
	}

	if len(headers) != 2 {
		t.Error("Number of headers is not as expected: %v", headers)
	}
}

func TestHeaderChangerAppend(t *testing.T) {
	headers := http.Header{}
	headers.Add("Already-Present", "1")

	changer := HeaderChanger{}

	changer["+Added"] = "Asdf"
	changer.changeHeaders(headers)
	if headers.Get("Added") != "Asdf" {
		t.Error("Missing header wasn't added")
	}

	changer["+Already-Present"] = "2"
	changer.changeHeaders(headers)
	if strings.Join(headers["Already-Present"], ", ") != "1, 2" {
		fmt.Println(headers)
		t.Errorf("Existing header wasn't appended: %v", headers["Already-Present"])
	}

	if len(headers) != 2 {
		t.Error("Number of headers is not as expected")
	}
}

func TestHeaderChangerDelete(t *testing.T) {
	headers := http.Header{}
	headers.Add("Already-Present", "1")

	changer := HeaderChanger{}

	changer["-Added"] = "Asdf"
	changer.changeHeaders(headers)
	if headers.Get("Added") != "" {
		t.Error("Missing header wasn't deleted")
	}

	changer["-Already-Present"] = "2"
	changer.changeHeaders(headers)
	if strings.Join(headers["Already-Present"], ", ") != "" {
		fmt.Println(headers)
		t.Errorf("Existing header wasn't deleted: %v", headers["Already-Present"])
	}

	if len(headers) != 0 {
		t.Error("Number of headers is not as expected")
	}
}
