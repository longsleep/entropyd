// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Url       string `json:"url"`
	Key       string `json:"key"`
	Watermark int    `json:"watermark"`
}

func LoadConfig(fn string) (*Config, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)

	if config.Url == "" {
		config.Url = "127.0.01:3344"
	}

	return config, err
}
