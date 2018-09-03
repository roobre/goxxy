package modules

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"roob.re/goxxy/tests"
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
		t.Errorf("Number of headers is not as expected: %v", headers)
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

func TestHeaderChangerRequest(t *testing.T) {
	request := tests.Get()
	request.Header.Add("Removing", "Whatever")

	prevAppend := strings.Join(request.Header["User-Agent"], ", ")

	changer := HeaderChanger{}
	changer["-Date"] = ""
	changer["+User-Agent"] = "mangled"
	changer["New"] = "NewHeader"
	handler := changer.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Date") != "" {
			t.Error("Date header wasnt removed")
		}

		if strings.Join(r.Header["User-Agent"], ", ") != prevAppend+", mangled" {
			t.Errorf("Error appending existing header: %v", r.Header)
		}

		if r.Header.Get("New") != "NewHeader" {
			t.Error("New header wasn't set")
		}
	}))

	handler.ServeHTTP(httptest.NewRecorder(), request)
}

func TestHeaderChangerResponse(t *testing.T) {
	response := tests.GetResponse()
	response.Header.Add("Removing", "Whatever")

	prevAppend := strings.Join(response.Header["Server"], ", ")

	changer := HeaderChanger{}
	changer["-Date"] = ""
	changer["+Server"] = "mangled"
	changer["New"] = "NewHeader"
	newresponse := changer.Mangle(response)

	if newresponse.Header.Get("Date") != "" {
		t.Error("Date header wasnt removed")
	}

	if strings.Join(response.Header["Server"], ", ") != prevAppend+", mangled" {
		t.Error("Error appending existing header")
	}

	if newresponse.Header.Get("New") != "NewHeader" {
		t.Error("New header wasn't set")
	}
}
