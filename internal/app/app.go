package app

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/mlc911/mlcsshtrayagent/internal/monitor"
)

//go:embed assets/icon.ico
var trayIcon []byte

//go:embed assets/icon.png
var notificationIcon []byte

const pollInterval = 10 * time.Second

type App struct {
	provider monitor.Provider
	done     chan struct{}
	labels   labels
}

func New(provider monitor.Provider) *App {
	return &App{provider: provider, done: make(chan struct{}), labels: newLabels()}
}

func (a *App) Run() {
	release, acquired, err := acquireSingleInstance()
	if err != nil || !acquired {
		return
	}
	defer release()

	enableNativeDarkMode()
	systray.Run(a.onReady, func() {})
}

func (a *App) onReady() {
	beeep.AppName = "MLC SSH Tray Agent"
	systray.SetIcon(trayIcon)
	systray.SetTitle("SSH/RDP")
	systray.SetTooltip(a.labels.tooltip)

	status := systray.AddMenuItem(a.labels.search, a.labels.searchTooltip)
	status.Disable()
	sessionItems := make([]*systray.MenuItem, 4)
	sessionItems[0] = systray.AddMenuItem(a.labels.empty, a.labels.emptyTooltip)
	sessionItems[0].Disable()
	for index := 1; index < len(sessionItems); index++ {
		sessionItems[index] = systray.AddMenuItem("", "")
		sessionItems[index].Hide()
	}
	systray.AddSeparator()
	aboutTitle, aboutTooltip, aboutMessage := aboutContent(a.labels)
	about := systray.AddMenuItem(aboutTitle, aboutTooltip)
	go func() {
		<-about.ClickedCh
		showAbout(aboutTitle, aboutMessage)
	}()
	systray.AddSeparator()
	quit := systray.AddMenuItem(a.labels.quit, a.labels.quitTooltip)

	go a.monitor(status, sessionItems)
	go func() {
		<-quit.ClickedCh
		close(a.done)
		systray.Quit()
	}()
}

func (a *App) monitor(status *systray.MenuItem, sessionItems []*systray.MenuItem) {
	var previous monitor.Snapshot
	first := true
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		snapshot, err := a.provider.Snapshot()
		if err != nil {
			status.SetTitle("Fehler: " + err.Error())
		} else {
			status.SetTitle(formatSnapshot(snapshot, a.labels))
			updateSessionItems(sessionItems, snapshot, a.labels)
			if !first {
				if err := notifyNewSessions(previous, snapshot, a.labels); err != nil {
					status.SetTitle(a.labels.notificationFailed)
				}
			}
			previous = snapshot
			first = false
		}

		select {
		case <-ticker.C:
		case <-a.done:
			return
		}
	}
}

func formatSnapshot(snapshot monitor.Snapshot, text labels) string {
	if len(snapshot.Sessions) == 0 {
		return text.empty
	}
	return fmt.Sprintf("%d %s", len(snapshot.Sessions), text.activeSessions)
}

func updateSessionItems(items []*systray.MenuItem, snapshot monitor.Snapshot, text labels) {
	if len(snapshot.Sessions) == 0 {
		items[0].SetTitle(text.empty)
		items[0].Disable()
		items[0].Show()
	}
	for index, item := range items {
		if index >= len(snapshot.Sessions) {
			if index > 0 {
				item.SetTitle("")
				item.Hide()
			}
			continue
		}
		session := snapshot.Sessions[index]
		title := string(session.Kind) + ": " + session.User
		if session.Source != "" {
			title += " (" + session.Source + ")"
		}
		if !session.StartedAt.IsZero() {
			title += " [" + formatDuration(time.Since(session.StartedAt)) + "]"
		}
		item.SetTitle(title)
		item.Show()
		item.Enable()
	}
}

func formatDuration(duration time.Duration) string {
	if duration < 0 {
		duration = 0
	}
	if duration < time.Minute {
		return fmt.Sprintf("%ds", int(duration/time.Second))
	}
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration/time.Minute))
	}
	return fmt.Sprintf("%dh %02dm", int(duration/time.Hour), int(duration/time.Minute)%60)
}

func notifyNewSessions(previous, current monitor.Snapshot, text labels) error {
	known := make(map[string]struct{}, len(previous.Sessions))
	for _, session := range previous.Sessions {
		known[session.Key()] = struct{}{}
	}

	for _, session := range current.Sessions {
		if _, exists := known[session.Key()]; exists {
			continue
		}
		message := string(session.Kind) + ": " + session.User
		if session.Source != "" {
			message += " " + text.from + " " + session.Source
		}
		if err := beeep.Notify(text.newSession, message, notificationIcon); err != nil {
			return err
		}
	}
	return nil
}
