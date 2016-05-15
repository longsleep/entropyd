// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/longsleep/entropyd"
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

	client, err := entropyd.NewClient([]byte(config.Key))
	if err != nil {
		log.Fatal("failed to create client: ", err)
	}

	entropy := entropyd.NewEntropy()
	if config.Watermark > 0 {
		err = entropy.SetWatermark(config.Watermark)
		if err != nil {
			log.Fatal("failed to set watermark: ", err)
		}
	}

	api := &API{config, client, entropy}

	http.Handle("/entropy/urandom", &APIHandler{api, UrandomHandler})

	wg := sync.WaitGroup{}
	for _, listener := range config.Listen {
		wg.Add(1)
		go func(l string) {
			log.Println("listening on", listener)
			err := http.ListenAndServe(l, nil)
			if err != nil {
				log.Fatal("failed to listen: ", err)
			}
		}(listener)
	}

	wg.Wait()
}

func usage() {
	fmt.Printf("Usage: %s <path-to-config>\n", os.Args[0])
}
