// Copyright 2016 Simon Eisenmann <simon@longsleep.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package random

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/longsleep/entropyd/fdset"
)

type Random struct {
	fp *os.File
}

func NewRandom() (*Random, error) {
	fp, err := os.OpenFile(randomDevice, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &Random{fp}, nil
}

func (random *Random) GetEntCnt() (int, error) {
	var cnt int

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, random.fp.Fd(), RNDGETENTCNT, uintptr(unsafe.Pointer(&cnt)))
	if errno != 0 {
		return 0, fmt.Errorf("failed to get current entropy level from kernel: %v", errno)
	}

	return cnt, nil
}

func (random *Random) Wait() error {
	wfds := &fdset.FdSet{}
	for {
		wfds.Zero()
		wfds.Set(random.fp.Fd())
		n, err := syscall.Select(int(random.fp.Fd())+1, nil, (*syscall.FdSet)(wfds), nil, nil)
		if err != nil {
			return err
		}
		if n >= 0 {
			break
		}
	}

	return nil
}

func (random *Random) Write(b []byte, entropyPerByte int) error {
	length := len(b)
	output := &RandPoolInfo{length * entropyPerByte, length, uintptr(unsafe.Pointer(&b))}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, random.fp.Fd(), RNDADDENTROPY, uintptr(unsafe.Pointer(output)))
	if errno != 0 {
		return fmt.Errorf("failed to add entropy to kernel: %v", errno)
	}

	return nil
}
