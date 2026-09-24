//go:build !windows

package monitor

// NewProvider returns a harmless provider for development and tests on non-Windows systems.
func NewProvider() Provider {
	return unsupportedProvider{}
}

type unsupportedProvider struct{}

func (unsupportedProvider) Snapshot() (Snapshot, error) {
	return Snapshot{}, nil
}
