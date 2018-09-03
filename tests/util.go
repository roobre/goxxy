package tests

// util.go contains utilities to generate mock requests and responses

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const RequestURL = "http://www.example.org/items"
const RequestPostdata = "testing=1&example=1&valid=0&testString=hewwo+world"
const RequestUA = "Test mocker"

const ResponseServer = "mock test server"
const ResponseHTML = `
<html>
<head>
	<title>Sample webpage</title>
</head>
<body>
	<h1>Main header</h1>
	<p>Paragrapgh describing header, click
	<a href="http://www.example.org/details">here</a> 
	<a href="https://www.example.org/details">(secure version)</a>
	for more info</p>
</body>
</html>
`

var ResponseJSON = struct {
	Name    string
	Count   int
	Value   string
	Ratio   float32
	Details struct{ Detail1 []string }
}{
	"Sample JSON",
	3,
	"ComplexValue",
	3.14,
	struct{ Detail1 []string }{[]string{"Very Detailed", "Yes yes"}},
}

func Get() *http.Request {
	req, err := http.NewRequest(http.MethodGet, RequestURL, bytes.NewReader(nil))
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", RequestUA)

	return req
}

func Post() *http.Request {
	req, err := http.NewRequest(http.MethodPost, RequestURL, strings.NewReader(RequestPostdata))
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", RequestUA)

	return req
}

func PostJson() *http.Request {
	r := Post()
	buf, err := json.Marshal(ResponseJSON)
	if err != nil {
		panic(err)
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(buf))
	r.Header.Set("Content-Type", "application/json")
	return r
}

func ResponseBoilerplate() *http.Response {
	h := http.Header{}
	h.Set("Server", ResponseServer)
	h.Set("Date", time.Now().Format(time.RFC1123))
	return &http.Response{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Header:     h,
	}
}

func GetResponse() *http.Response {
	resp := ResponseBoilerplate()
	resp.Request = Get()
	resp.ContentLength = int64(len(ResponseHTML))
	resp.Body = ioutil.NopCloser(strings.NewReader(ResponseHTML))

	return resp
}

func GetResponseJSON() *http.Response {
	resp := ResponseBoilerplate()
	resp.Request = Get()

	buf, err := json.Marshal(ResponseJSON)
	if err != nil {
		panic(err)
	}
	resp.ContentLength = int64(len(buf))
	resp.Body = ioutil.NopCloser(bytes.NewReader(buf))

	resp.Header.Set("Content-Type", "application/json")

	return resp
}
