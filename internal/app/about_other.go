//go:build !windows

package app

import "fmt"

func showAbout(title, message string) {
	fmt.Printf("%s\n%s\n", title, message)
}
