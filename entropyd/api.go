// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/longsleep/entropyd"
)

type API struct {
	Config  *Config
	Client  *entropyd.Client
	Entropy *entropyd.Entropy
}

type APIError interface {
	error
	Status() int
}

type APIStatusError struct {
	Code int
	Err  error
}

func (err APIStatusError) Error() string {
	return err.Err.Error()
}

func (err APIStatusError) Status() int {
	return err.Code
}
