//go:build !windows

package baza

import (
	"errors"
	"syscall"
)

// luaBlocareSO cere blocarea exclusivă a sistemului de operare pe fișier, fără
// să aștepte. Sistemul o eliberează singur când procesul se termină, oricum.
// „false” fără eroare înseamnă că o ține altcineva.
func luaBlocareSO(fd uintptr) (bool, error) {
	for {
		err := syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
		switch {
		case err == nil:
			return true, nil
		case errors.Is(err, syscall.EINTR):
			continue
		case errors.Is(err, syscall.EWOULDBLOCK):
			return false, nil
		default:
			return false, err
		}
	}
}
