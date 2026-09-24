package monitor

import "time"

// SessionKind identifies the operating-system service that owns a session.
type SessionKind string

const (
	SessionKindSSH SessionKind = "SSH"
	SessionKindRDP SessionKind = "RDP"
	SessionKindWSL SessionKind = "WSL"
)

// Session describes one currently active remote or subsystem session.
type Session struct {
	Kind      SessionKind
	User      string
	Source    string
	StartedAt time.Time
	ID        string
}

// Snapshot is the monitor state at one point in time.
type Snapshot struct {
	Sessions []Session
	At       time.Time
}

// Count returns the number of sessions in the snapshot.
func (s Snapshot) Count() int {
	return len(s.Sessions)
}

// Key returns a stable identity for change detection.
func (s Session) Key() string {
	return string(s.Kind) + ":" + s.ID + ":" + s.User
}
