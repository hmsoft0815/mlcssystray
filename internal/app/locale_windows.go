//go:build windows

package app

import (
	"syscall"
	"unsafe"
)

func systemLanguage() string {
	buffer := make([]uint16, 85)
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")
	if result, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer))); result != 0 {
		return syscall.UTF16ToString(buffer)
	}
	return ""
}
