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
	defaultNRJMXExec = "c:\\progra~1\\newrel~1\\nrjmx\\nrjmx.bat"
)

// buildExecCommand adds os specifics to the command.
func buildExecCommand(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, filepath.Clean(getNRJMXExec()), nrJMXV2Flag)
}
