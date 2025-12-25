package windows

import (
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

	procGetClientRect = user32.NewProc("GetClientRect")
	procFindWindowExW = user32.NewProc("FindWindowExW")
	procMoveWindow    = user32.NewProc("MoveWindow")

	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procBringWindowToTop    = user32.NewProc("BringWindowToTop")
	procShowWindow          = user32.NewProc("ShowWindow")
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
)

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
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

func BringWindowToTop(hwnd uintptr) bool {
	ret, _, _ := procBringWindowToTop.Call(hwnd)
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

func HideTitleBar(hwnd uintptr) {
	MakeFramelessAndFitWebView(hwnd)
}

func hideTitleBar(hwnd uintptr) {
	style := getStyle(hwnd)

	style &^= WS_CAPTION
	style &^= WS_SYSMENU
	style &^= WS_THICKFRAME
	style &^= WS_MINIMIZEBOX
	style &^= WS_MAXIMIZEBOX

	setStyle(hwnd, style)
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
