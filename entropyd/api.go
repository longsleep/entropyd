// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/longsleep/entropyd"
)

type API struct {
	Config  *Config
	Client  *entropyd.Client
	Entropy *entropyd.Entropy

	allowNet []*net.IPNet
	allowIP  []*net.IP
}

func NewAPI(config *Config, client *entropyd.Client, entropy *entropyd.Entropy) *API {
	api := &API{
		Config:  config,
		Client:  client,
		Entropy: entropy,
	}

	for _, addr := range config.Allow {
		if strings.Contains(addr, "/") {
			if _, cidr, err := net.ParseCIDR(addr); err == nil {
				api.allowNet = append(api.allowNet, cidr)
			}
		} else {
			if ip := net.ParseIP(addr); ip != nil {
				api.allowIP = append(api.allowIP, &ip)
			}
		}
	}

	return api
}

func (api *API) ValidateRequest(request *http.Request) error {
	remoteHost, _, _ := net.SplitHostPort(request.RemoteAddr)
	remoteIP := net.ParseIP(remoteHost)

	for _, cidr := range api.allowNet {
		if cidr.Contains(remoteIP) {
			return nil
		}
	}
	for _, ip := range api.allowIP {
		if ip.Equal(remoteIP) {
			return nil
		}
	}

	return errors.New("remote address not allowed")
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
