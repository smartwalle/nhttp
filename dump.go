package http4go

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func CopyBody(req *http.Request) (io.ReadCloser, error) {
	var err error
	if req.Body == nil {
		return http.NoBody, nil
	}

	r1, r2, err := drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = r1

	return r2, nil
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
