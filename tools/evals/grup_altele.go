//go:build !unix

package evals

import "os/exec"

func grupNou(cmd *exec.Cmd) {}
