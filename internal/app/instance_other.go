//go:build !windows

package app

func acquireSingleInstance() (func(), bool, error) {
	return func() {}, true, nil
}
