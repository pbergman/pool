package pool

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ResponseInterface implements the basic methods of http.response
type ResponseInterface interface {
	Cookies() []*http.Cookie
	Location() (*url.URL, error)
	ProtoAtLeast(major, minor int) bool
	Write(w io.Writer) error
}

// Response embedded response wrapper
type Response struct {
	*http.Response
}

// GetBody checks the content encoding of response and assign the appropriated reader to the body
func (r *Response) GetBody() (io.ReadCloser, error) {
	// Check that the server actual sent compressed data
	switch r.Header.Get("Content-Encoding") {
	case "deflate":
		reader := flate.NewReader(r.Body)
		return reader, nil
	case "gzip":
		reader, _ := gzip.NewReader(r.Body)
		return reader, nil
	default:
		return r.Body, nil
	}
}

// GetBodyString will process the body and return a string
func (r *Response) GetBodyBytes() ([]byte, error) {
	reader, err := r.GetBody()
	defer reader.Close()
	if err != nil {
		return make([]byte, 0), err
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return make([]byte, 0), err
	}
	return body, nil
}
