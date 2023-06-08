package nhttp

import (
	"bytes"
	"io"
	"net/http"
)

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a fork of net/http/httputil

func DumpBody(req *http.Request) (io.ReadCloser, error) {
	var err error
	if req == nil || req.Body == nil {
		return http.NoBody, nil
	}

	r1, r2, err := drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = r1

	return r2, nil
}

func drainBody(src io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if src == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(src); err != nil {
		return nil, nil, err
	}
	if err = src.Close(); err != nil {
		return nil, nil, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
