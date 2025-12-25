package windows

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32  = windows.NewLazySystemDLL("user32.dll")
	shell32 = windows.NewLazySystemDLL("shell32.dll")

	procRegisterHotKey      = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey    = user32.NewProc("UnregisterHotKey")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procShellExecuteW       = shell32.NewProc("ShellExecuteW")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procFindWindowW         = user32.NewProc("FindWindowW")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")

	procGetWindowLongPtrW = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW = user32.NewProc("SetWindowLongPtrW")

	procGetClientRect    = user32.NewProc("GetClientRect")
	procFindWindowExW    = user32.NewProc("FindWindowExW")
	procMoveWindow       = user32.NewProc("MoveWindow")
	procBringWindowToTop = user32.NewProc("BringWindowToTop")
	procShowWindow       = user32.NewProc("ShowWindow")

	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")

	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

const (
	GWL_STYLE = int32(-16)

	WS_CAPTION     = 0x00C00000
	WS_SYSMENU     = 0x00080000
	WS_THICKFRAME  = 0x00040000
	WS_MINIMIZEBOX = 0x00020000
	WS_MAXIMIZEBOX = 0x00010000

	SWP_NOMOVE       = 0x0002
	SWP_NOSIZE       = 0x0001
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020
	SWP_SHOWWINDOW   = 0x0040
)

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

func SetupAppWindow() {
	pid := uint32(os.Getpid())

	hwnd := GetWindowHandleByPID(pid)
	if hwnd == 0 {
		return
	}
	SetWindowTopMost(hwnd)
}

func SetWindowTopMost(hwnd uintptr) {
	hwndTopMost := ^uintptr(0)

	procSetWindowPos.Call(
		hwnd,
		hwndTopMost,
		0, 0, 0, 0,
		uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_SHOWWINDOW),
	)

	SetForegroundWindow(hwnd)
	procBringWindowToTop.Call(hwnd)
}

func GetWindowHandleByPID(targetPID uint32) uintptr {
	var foundHwnd uintptr = 0

	cb := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		var wndPid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&wndPid)))

		if wndPid == targetPID {
			isVisible, _, _ := procIsWindowVisible.Call(hwnd)
			if isVisible != 0 {
				foundHwnd = hwnd
				return 0
			}
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)
	return foundHwnd
}

func RegisterHotKey(hwnd uintptr, id int, mod uint32, vk uint32) bool {
	r, _, _ := procRegisterHotKey.Call(hwnd, uintptr(id), uintptr(mod), uintptr(vk))
	return r != 0
}

func UnregisterHotKey(hwnd uintptr, id int) {
	procUnregisterHotKey.Call(hwnd, uintptr(id))
}

func GetMessage(msg *MSG) int32 {
	ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0)
	return int32(ret)
}

func TranslateMessage(msg *MSG) {
	procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func DispatchMessage(msg *MSG) {
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func ShellOpen(path string) {
	verb := utf16Ptr("open")
	target := utf16Ptr(path)

	procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(target)),
		0,
		0,
		5,
	)
}

func GetForegroundWindow() uintptr {
	ret, _, _ := procGetForegroundWindow.Call()
	return ret
}

func SetForegroundWindow(hwnd uintptr) bool {
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	return ret != 0
}

func ShowWindow(hwnd uintptr, nCmdShow int) bool {
	ret, _, _ := procShowWindow.Call(hwnd, uintptr(nCmdShow))
	return ret != 0
}

func MakeFramelessAndFitWebView(hwnd uintptr) {
	hideTitleBar(hwnd)
	applyFrameChanged(hwnd)
	fitWebViewChildToClient(hwnd)
}

func hideTitleBar(hwnd uintptr) {
	style := getStyle(hwnd)

	// Quitamos CAPTION (la barra azul/blanca con el título)
	style &^= WS_CAPTION

	// Quitamos SYSMENU (los botones cerrar/min/max)
	style &^= WS_SYSMENU

	// ¡IMPORTANTE! Quitamos THICKFRAME.
	// Este es el culpable de ese marco blanco grueso que ves en la foto.
	style &^= WS_THICKFRAME

	// Quitamos cualquier otro borde estándar
	style &^= 0x00800000 // WS_BORDER

	style &^= WS_MINIMIZEBOX
	style &^= WS_MAXIMIZEBOX

	setStyle(hwnd, style)

	// Opcional: Agregar un estilo extendido para que tenga sombra (si quieres)
	// o dejarlo plano. Por ahora, plano para arreglar el error visual.
}

func applyFrameChanged(hwnd uintptr) {
	procSetWindowPos.Call(
		hwnd,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED),
	)
}

func fitWebViewChildToClient(hwnd uintptr) {
	var rc RECT
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))

	w := int32(rc.Right - rc.Left)
	h := int32(rc.Bottom - rc.Top)
	if w <= 0 || h <= 0 {
		return
	}

	child := findChildByClass(hwnd, "Chrome_WidgetWin_0")
	if child == 0 {
		child = findFirstChild(hwnd)
	}
	if child == 0 {
		return
	}

	procMoveWindow.Call(
		child,
		0, 0,
		uintptr(w),
		uintptr(h),
		1,
	)
}

func findChildByClass(parent uintptr, class string) uintptr {
	cls := utf16Ptr(class)
	h, _, _ := procFindWindowExW.Call(
		parent,
		0,
		uintptr(unsafe.Pointer(cls)),
		0,
	)
	return h
}

func findFirstChild(parent uintptr) uintptr {
	h, _, _ := procFindWindowExW.Call(parent, 0, 0, 0)
	return h
}

func getStyle(hwnd uintptr) uintptr {
	gwlStyle := int32(GWL_STYLE)
	s, _, _ := procGetWindowLongPtrW.Call(hwnd, uintptr(gwlStyle))
	return s
}

func setStyle(hwnd uintptr, style uintptr) {
	gwlStyle := int32(GWL_STYLE)
	procSetWindowLongPtrW.Call(hwnd, uintptr(gwlStyle), style)
}

func utf16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return windows.StringToUTF16Ptr("")
	}
	return p
}
