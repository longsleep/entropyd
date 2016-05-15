// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entropyd

import (
	"testing"
)

var (
	key []byte = []byte("AES256Key-32Characters1234567890")
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("exampleplaintext")

	client, err := NewClient(key)
	if err != nil {
		t.Error(err)
		return
	}

	nonce, err := client.Nonce()
	if err != nil {
		t.Error(err)
		return
	}

	ciphertext := client.Encrypt(nonce, plaintext)
	decryptedtext, err := client.Decrypt(nonce, ciphertext)
	if err != nil {
		t.Error(err)
		return
	}

	if string(decryptedtext) != string(plaintext) {
		t.Error("decrypted text does not match plain text")
		return
	}
}
