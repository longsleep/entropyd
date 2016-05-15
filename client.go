// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entropyd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Client struct {
	aead cipher.AEAD
	name string
}

func NewClient(key []byte) (*Client, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Client{aead, "aes-gcm"}, nil
}

func (c *Client) Encrypt(nonce []byte, plaintext []byte) []byte {
	return c.aead.Seal(nil, nonce, plaintext, nil)
}

func (c *Client) Decrypt(nonce []byte, ciphertext []byte) ([]byte, error) {
	return c.aead.Open(nil, nonce, ciphertext, nil)
}

func (c *Client) Nonce() ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func (c *Client) NonceSize() int {
	return c.aead.NonceSize()
}

func (c *Client) String() string {
	return c.name
}
