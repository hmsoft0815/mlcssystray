//go:build !windows

package app

import "os"

func systemLanguage() string {
	if language := os.Getenv("LC_ALL"); language != "" {
		return language
	}
	return os.Getenv("LANG")
}
