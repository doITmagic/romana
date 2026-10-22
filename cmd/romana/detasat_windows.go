//go:build windows

package main

import "syscall"

// detachedProcess: procesul nu primește consola părintelui.
const detachedProcess = 0x00000008

// atributeDetasat: fără consolă și în alt grup de procese, ca Ctrl-C să nu
// oprească actualizarea pornită în fundal.
func atributeDetasat() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: detachedProcess | syscall.CREATE_NEW_PROCESS_GROUP, HideWindow: true}
}
