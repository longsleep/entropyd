// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entropyd

import (
	"crypto/rand"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	PoolSizeProcFs   = "/proc/sys/kernel/random/poolsize"
	BaseContentType  = "application/x-entropy"
	NonceContentType = "application/x-entropy-nonce"
)

type Entropy struct {
}

func NewEntropy() *Entropy {
	return &Entropy{}
}

func (entropy *Entropy) Urandom(b []byte) error {
	_, err := rand.Read(b)

	return err
}

func (entropy *Entropy) GetPoolSize() (int, error) {
	poolSizeString, err := ioutil.ReadFile(PoolSizeProcFs)
	if err != nil {
		return 0, err
	}
	poolSize, err := strconv.Atoi(strings.TrimSpace(string(poolSizeString)))
	if err != nil {
		return 0, err
	}

	return poolSize, nil
}
