package monitor

// Provider discovers sessions from one operating-system source.
type Provider interface {
	Snapshot() (Snapshot, error)
}
