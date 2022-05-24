/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"os/exec"
	"path/filepath"
)

const (
	// defaultNRJMXExec default nrjmx tool executable path.
	defaultNRJMXExec = "/usr/local/bin/nrjmx"
)

// buildExecCommand adds os specifics to the command.
func buildExecCommand(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, filepath.Clean(getNRJMXExec()), nrJMXV2Flag)
}
