// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entropyd

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

type HTTPClient struct {
	*Client
	uri string
}

func NewHTTPClient(key []byte, uri string) (*HTTPClient, error) {
	client, err := NewClient(key)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{client, uri}, nil
}

func (c *HTTPClient) NewUrandomRequest(length int, nonce []byte) (*http.Request, error) {
	uri := fmt.Sprintf("%s/entropy/urandom", c.uri)

	req, err := http.NewRequest("POST", uri, bytes.NewReader(nonce))
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Set("length", strconv.Itoa(length))
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", NonceContentType)

	return req, nil
}
