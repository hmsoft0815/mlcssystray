package app

import "golang.org/x/sys/windows"

func acquireSingleInstance() (func(), bool, error) {
	name, err := windows.UTF16PtrFromString(`Local\MLCSSHTrayAgent`)
	if err != nil {
		return nil, false, err
	}

	handle, err := windows.CreateMutex(nil, false, name)
	if err == windows.ERROR_ALREADY_EXISTS {
		_ = windows.CloseHandle(handle)
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return func() {
		_ = windows.CloseHandle(handle)
	}, true, nil
}
