// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/longsleep/entropyd"
	"github.com/longsleep/entropyd/random"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	config, err := LoadConfig(args[0])
	if err != nil {
		log.Fatal("failed to parse config: ", err)
	}

	// Prepare entropyd client.
	client, err := entropyd.NewHTTPClient([]byte(config.Key), config.Url)
	if err != nil {
		log.Fatal("failed to create client: ", err)
	}

	// Open random device.
	random, err := random.NewRandom()
	if err != nil {
		log.Fatal("failed to open random device: ", err)
	}

	// Get Kernel entropy pool size.
	entropy := entropyd.NewEntropy()
	poolSize, err := entropy.GetPoolSize()
	if err != nil {
		log.Fatal("failed to get entropy pool size: ", err)
	}

	// Prepare http client for connection reuse.
	conn := &http.Client{}

	var nonce []byte
	var current int
	var nbytes int

	var req *http.Request
	var resp *http.Response
	var body []byte
	var plaintext []byte
	var contentType string
	var responseType string = fmt.Sprintf("%s-%s", entropyd.BaseContentType, client)

	log.Printf("entropyd URL: %s\n", config.Url)
	for {
		if err != nil {
			// Wait a bit if there was an error before.
			time.Sleep(1 * time.Second)
		}

		// Generate random nonce for AES-GCM.
		nonce, err = client.Nonce()
		if err != nil {
			log.Fatal("failed to create nonce: ", err)
		}

		// Wait for random device to be ready for write.
		err = random.Wait()
		if err != nil {
			log.Fatal("failed to wait for random device: ", err)
		}

		// Get current entropy in bits.
		current, err = random.GetEntCnt()
		if err != nil {
			log.Fatal("failed to get entropy count: ", err)
		}
		// Compute needed entropy bytes.
		nbytes = (poolSize - current) / 8
		if nbytes < 1 {
			continue
		}
		// Server can do 2048 byte max.
		if nbytes > 2048 {
			nbytes = 2048
		}

		// Prepare request.
		req, err = client.NewUrandomRequest(nbytes, nonce)
		if err != nil {
			log.Fatal("failed to create request: ", err)
		}
		// Send request.
		resp, err = conn.Do(req)
		if err != nil {
			log.Println("request failed:", err)
			continue
		}

		// Read response.
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// Error handling.
		if resp.StatusCode != 200 {
			err = fmt.Errorf("invalid response: %v", string(resp.Status))
			log.Println(err)
			continue
		}

		// Parse response.
		contentType = resp.Header.Get("Content-Type")
		switch contentType {
		case responseType:
			// AES-GCM response decrypt/validate.
			plaintext, err = client.Decrypt(nonce, body)
			if err != nil {
				log.Println("failed to decrypt:", err)
				continue
			}
			// Write decrypted plaintext to random device.
			err = random.Write(plaintext, 8)
			if err != nil {
				log.Println("failed to write:", err)
				continue
			}
		default:
			err = fmt.Errorf("invalid response type: %v", contentType)
			log.Println(err)
			continue
		}
	}
}

func usage() {
	fmt.Printf("Usage: %s <path-to-config>\n", os.Args[0])
}
