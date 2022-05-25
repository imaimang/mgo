//go:build !windows
// +build !windows

package console

import (
	"fmt"
)

func ColorPrint(messageType string, message string) {
	var color string
	switch messageType {
	case "DEBUG":
		color = "34"
	case "INFO":
		color = "32"
	case "WARN":
		color = "33"
	case "ERROR":
		color = "31"
	}
	fmt.Printf("\033[1;"+color+"m%s\033[0m\r\n", message)
}
