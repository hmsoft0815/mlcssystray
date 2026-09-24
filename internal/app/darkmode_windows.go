//go:build windows

package app

import "syscall"

func enableNativeDarkMode() {
	proc := syscall.NewLazyDLL("uxtheme.dll").NewProc("SetPreferredAppMode")
	if err := proc.Find(); err != nil {
		return
	}
	const allowDark = 2
	proc.Call(uintptr(allowDark))
}
