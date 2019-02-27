package cdata

import (
	"golang.org/x/sys/windows"
)

func Ref() int64 {
	return int64(windows.GetCurrentThreadId())
}
