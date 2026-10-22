//go:build windows

package baza

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	lockfileFailImmediately = 0x1
	lockfileExclusiveLock   = 0x2
	errorLockViolation      = syscall.Errno(33)
)

var procLockFileEx = syscall.NewLazyDLL("kernel32.dll").NewProc("LockFileEx")

// luaBlocareSO cere blocarea exclusivă a sistemului de operare pe fișier, fără
// să aștepte (LockFileEx). Sistemul o eliberează singur când procesul se
// termină. „false” fără eroare înseamnă că o ține altcineva.
func luaBlocareSO(fd uintptr) (bool, error) {
	var ov syscall.Overlapped
	r, _, err := procLockFileEx.Call(fd, lockfileExclusiveLock|lockfileFailImmediately,
		0, 1, 0, uintptr(unsafe.Pointer(&ov)))
	if r != 0 {
		return true, nil
	}
	if errors.Is(err, errorLockViolation) {
		return false, nil
	}
	return false, err
}
