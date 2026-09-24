package app

import (
	"os"
	"strings"
)

type labels struct {
	tooltip            string
	search             string
	searchTooltip      string
	empty              string
	emptyTooltip       string
	quit               string
	quitTooltip        string
	notificationFailed string
	newSession         string
	from               string
	activeSessions     string
}

func newLabels() labels {
	language := strings.ToLower(os.Getenv("MLCSSHTRAY_LANG"))
	if language == "" {
		language = strings.ToLower(systemLanguage())
	}
	if strings.HasPrefix(language, "de") {
		return labels{
			tooltip:            "Aktive SSH- und RDP-Sitzungen",
			search:             "Suche Sitzungen ...",
			searchTooltip:      "Aktive Sitzungen",
			empty:              "Keine SSH/RDP-Sitzungen",
			emptyTooltip:       "Aktive SSH- und RDP-Sitzungen",
			quit:               "Beenden",
			quitTooltip:        "Tray-Agent beenden",
			notificationFailed: "Benachrichtigung fehlgeschlagen",
			newSession:         "Neue Sitzung",
			from:               "von",
			activeSessions:     "aktive Sitzung(en)",
		}
	}
	return labels{
		tooltip:            "Active SSH and RDP sessions",
		search:             "Searching sessions ...",
		searchTooltip:      "Active sessions",
		empty:              "No SSH/RDP sessions",
		emptyTooltip:       "Active SSH and RDP sessions",
		quit:               "Quit",
		quitTooltip:        "Quit tray agent",
		notificationFailed: "Notification failed",
		newSession:         "New session",
		from:               "from",
		activeSessions:     "active session(s)",
	}
}
