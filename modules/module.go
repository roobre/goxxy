package modules // import "roob.re/goxxy/modules"

import (
	"bytes"
	"io"
	"net/http"
)

const defaultResponseBufferSize = 1024
const defaultResponseMaxSize = 1024 * 1024 * 1024

type maxSizer struct {
	MaxSize int64
}

func (m *maxSizer) maxSize() int64 {
	if m.MaxSize != 0 {
		return m.MaxSize
	} else {
		return defaultResponseMaxSize
	}
}

type byter interface {
	Bytes() []byte
	UnreadByte() error
}

type bufferCloser struct {
	*bytes.Buffer
}

func (b bufferCloser) Close() error {
	b.Buffer.Reset()
	return nil
}

// CopyBody reads the whole response body into a io.Buffer, and returns the slice of bytes from it as well as the reader buffer. It also sets the response body to the new Buffer.
// Warning: changes to the returned byte slice may not be reflected into the response automatically, if it is resliced somewhere. If you're unsure, re-set response.Body to a new buffer from the slice again.
func CopyBody(response *http.Response) (body []byte) {
	// If we already did the copy (response.Body implements `Bytes()`) and the buffer is unread (UnreadByte returns non-nil), just return those bytes
	if byter, isBuffer := response.Body.(byter); isBuffer && byter.UnreadByte() != nil {
		return byter.Bytes()
	}

	responseLen := response.ContentLength
	if responseLen < 0 {
		responseLen = defaultResponseBufferSize
	}
	buffer := bytes.NewBuffer(make([]byte, 0, responseLen))
	io.Copy(buffer, response.Body)
	response.Body.Close()

	response.Body = bufferCloser{buffer}
	return buffer.Bytes()
}
