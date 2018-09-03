package modules

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"roob.re/goxxy/tests"
	"strings"
	"testing"
)

func TestCopyBody(t *testing.T) {
	constBuffer := ioutil.NopCloser(strings.NewReader(tests.ResponseHTML))
	resp := http.Response{Body: constBuffer}

	//	First call to CopyBody should assign a new buffer
	body := CopyBody(&resp)

	if bytes.Compare(body, []byte(tests.ResponseHTML)) != 0 {
		t.Error("Returned and original body are different")
	}

	if resp.Body == constBuffer {
		t.Error("Original Response.Body was not modified")
	}

	oldBody := resp.Body

	// Call it a few times more
	for i := 0; i < 10; i++ {
		body = CopyBody(&resp)
		if bytes.Compare(body, []byte(tests.ResponseHTML)) != 0 {
			t.Error("Returned 2 and original body are different")
			t.FailNow()
		}
		if resp.Body != oldBody {
			t.Error("Body was unnecessarily re-copied")
			t.FailNow()
		}

		if _, isBuffer := resp.Body.(byter); !isBuffer {
			t.Error("Body is not a byter")
			t.FailNow()
		}
	}

	// Replace body and re-check
	oldBody = resp.Body
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	body = CopyBody(&resp)
	if bytes.Compare(body, []byte(tests.ResponseHTML)) != 0 {
		t.Error("Returned 2 and original body are different")
	}

	if resp.Body == oldBody {
		t.Error("Body was mistaken as a clean buffer and not recopied")
	}
}

// I'm doing this only because seeing coverage in green makes me happy
func TestMaxSizer(t *testing.T) {
	m := maxSizer{}

	if m.maxSize() != defaultResponseMaxSize {
		t.Fail()
	}

	m.MaxSize = 15
	if m.maxSize() != 15 {
		t.Fail()
	}
}
