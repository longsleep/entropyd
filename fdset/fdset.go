// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fdset

import (
	"syscall"
)

type FdSet syscall.FdSet

func (fds *FdSet) Set(fd uintptr) {
	fds.Bits[fd/NFDBITS] |= (1 << (fd % NFDBITS))
}

func (fds *FdSet) Zero() {
	for i := range fds.Bits {
		fds.Bits[i] = 0
	}
}
