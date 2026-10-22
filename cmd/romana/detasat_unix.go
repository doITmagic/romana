//go:build !windows

package main

import "syscall"

// atributeDetasat: sesiune nouă, ca Ctrl-C din terminal să nu oprească
// actualizarea pornită în fundal.
func atributeDetasat() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
