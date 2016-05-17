// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/longsleep/entropyd"
)

func UrandomHandler(api *API, w http.ResponseWriter, r *http.Request) error {
	err := api.ValidateRequest(r)
	if err != nil {
		http.Error(w, "address not allowed", http.StatusForbidden)
		return nil
	}

	r.ParseForm()

	length := 64
	if lengthString := r.Form.Get("length"); lengthString != "" {
		length, err = strconv.Atoi(lengthString)
	}

	if length < 1 || length > 2048 {
		http.Error(w, "invalid length", http.StatusBadRequest)
		return nil
	}

	var nonce []byte

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case entropyd.NonceContentType:
		nonce, err = ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}
		if len(nonce) != api.Client.NonceSize() {
			http.Error(w, "invalid nonce size", http.StatusBadRequest)
			return nil
		}
	default:
		http.Error(w, "invalid request type", http.StatusBadRequest)
		return nil
	}

	b := make([]byte, length)
	err = api.Entropy.Urandom(b)
	if err != nil {
		return err
	}

	c := api.Client.Encrypt(nonce, b)

	w.Header().Set("Content-Type", fmt.Sprintf("%s-%s", entropyd.BaseContentType, api.Client))
	w.Write(c)

	log.Printf("providing %d bytes to %v", length, r.RemoteAddr)

	return nil
}
