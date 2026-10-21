//go:build unix

package evals

import (
	"os/exec"
	"syscall"
)

// grupNou pune procesul într-un grup propriu și, la anulare, omoară tot grupul.
func grupNou(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
}
