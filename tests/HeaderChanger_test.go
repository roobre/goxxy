package tests

import (
	"net/http"
	"net/http/httptest"
	"roob.re/goxxy/modules"
	"strings"
	"testing"
)

func TestHeaderChangerRequest(t *testing.T) {
	request := Get()
	request.Header.Add("Removing", "Whatever")

	prevAppend := strings.Join(request.Header["User-Agent"], ", ")

	changer := modules.HeaderChanger{}
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
	response := GetResponse()
	response.Header.Add("Removing", "Whatever")

	prevAppend := strings.Join(response.Header["Server"], ", ")

	changer := modules.HeaderChanger{}
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
