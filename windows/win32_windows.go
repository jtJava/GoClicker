//go:build windows
// +build windows

package windows

import (
	"syscall"
	"unsafe"
)

var (
	modUser32, _               = syscall.LoadDLL("user32.dll")
	procPostMessageA, _        = modUser32.FindProc("PostMessageA")
	procGetForegroundWindow, _ = modUser32.FindProc("GetForegroundWindow")
	procGetClassNameW, _       = modUser32.FindProc("GetClassNameW")
	procGetAsyncKeyState, _    = modUser32.FindProc("GetAsyncKeyState")
)

func GetClassNameW(hwnd uintptr) string {
	buf := make([]uint16, 255)
	procGetClassNameW.Call(
		hwnd,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(255))

	return syscall.UTF16ToString(buf)
}

func GetForegroundWindow() uintptr {
	r, _, _ := procGetForegroundWindow.Call()

	return r
}

func PostMessage(hhk uintptr, code uintptr, wParam, lParam uintptr) uintptr {
	r, _, _ := procPostMessageA.Call(hhk, code, wParam, lParam)

	return r
}

func GetAsyncKeyState(vKey int) uint16 {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))

	return uint16(r)
}
