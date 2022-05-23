/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	// defaultNRJMXExec default nrjmx tool executable path.
	defaultNRJMXExec = "/usr/bin/nrjmx"
)

// buildExecCommand adds os specifics to the command.
func buildExecCommand(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, filepath.Clean(getNRJMXExec()), nrJMXV2Flag)

	// Terminate the subprocess when parent dies.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}

	return cmd
}
