package app

import "fmt"

func aboutContent(text labels) (title, tooltip, message string) {
	if text.newSession == "Neue Sitzung" {
		return "Über MLC SSH Tray Agent", "Informationen zum Tray-Agent", fmt.Sprintf("MLC SSH Tray Agent v%s\nCopyright Michael Lechner\nLizenziert unter MIT.", Version)
	}
	return "About MLC SSH Tray Agent", "About this tray agent", fmt.Sprintf("MLC SSH Tray Agent v%s\nCopyright Michael Lechner\nLicensed under MIT.", Version)
}