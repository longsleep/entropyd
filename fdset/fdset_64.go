// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!netbsd,!openbsd
// +build amd64 arm64

package fdset

const (
	NFDBITS = 8 * 8
)
