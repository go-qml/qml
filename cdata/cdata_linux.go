package cdata

import "syscall"

func Ref() int64 {
	return int64(syscall.Gettid())
}
