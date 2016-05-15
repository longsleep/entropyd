// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
)

type APIHandler struct {
	*API
	H func(a *API, w http.ResponseWriter, r *http.Request) error
}

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.API, w, r)
	if err != nil {
		switch e := err.(type) {
		case APIError:
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
