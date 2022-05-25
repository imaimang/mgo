//go:build windows
// +build windows

package console

import (
	"fmt"
	"syscall"
)

var (
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
)

func ColorPrint(messageType string, message string) {
	var color int
	switch messageType {
	case "DEBUG":
		color = 9
	case "INFO":
		color = 10
	case "WARN":
		color = 14
	case "ERROR":
		color = 12
	}
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(color))
	fmt.Println(message)
	CloseHandle.Call(handle)

	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(15))
	CloseHandle.Call(handle)
}
