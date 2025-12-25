package hotkey

import (
	"runtime"

	"winray-app/internal/windows"
)

const (
	ID_HOTKEY   = 1
	WM_HOTKEY   = 0x0312
	MOD_CONTROL = 0x0002
	VK_F        = 0x46
)

type ToggleFunc func()

func Loop(toggle ToggleFunc) {
	runtime.LockOSThread()

	if !windows.RegisterHotKey(0, ID_HOTKEY, MOD_CONTROL, VK_F) {
		panic("RegisterHotKey failed: Ctrl+F might be used by another app")
	}
	defer windows.UnregisterHotKey(0, ID_HOTKEY)

	var msg windows.MSG
	for {
		ret := windows.GetMessage(&msg)
		if ret <= 0 {
			break
		}

		if msg.Message == WM_HOTKEY && int(msg.WParam) == ID_HOTKEY {
			toggle()
			continue
		}

		windows.TranslateMessage(&msg)
		windows.DispatchMessage(&msg)
	}
}
