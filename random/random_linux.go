// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package random

// #include <linux/random.h>
import "C"

const (
	randomDevice = "/dev/random"

	RNDGETENTCNT  = C.RNDGETENTCNT
	RNDADDENTROPY = C.RNDADDENTROPY
)

type RandPoolInfo struct {
	entropy_count int
	buf_size      int
	buf           uintptr
}
