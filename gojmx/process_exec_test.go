/*
 * Copyright 2026 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSiblingNRJMXExec(t *testing.T) {
	t.Run("finds the tool next to the executable", func(t *testing.T) {
		dir := resolvedTempDir(t)
		executable := touch(t, dir, "integration")
		tool := touch(t, dir, nrJMXFileName)

		assert.Equal(t, tool, siblingNRJMXExec(executable))
	})

	t.Run("follows a symlink to the executable", func(t *testing.T) {
		realDir, linkDir := resolvedTempDir(t), resolvedTempDir(t)
		executable := touch(t, realDir, "integration")
		tool := touch(t, realDir, nrJMXFileName)
		link := filepath.Join(linkDir, "integration")
		if err := os.Symlink(executable, link); err != nil {
			t.Skipf("symlinks not supported: %v", err)
		}

		assert.Equal(t, tool, siblingNRJMXExec(link))
	})

	t.Run("handles directories with spaces", func(t *testing.T) {
		dir := filepath.Join(resolvedTempDir(t), "with space")
		require.NoError(t, os.Mkdir(dir, 0o755))
		executable := touch(t, dir, "integration")
		tool := touch(t, dir, nrJMXFileName)

		assert.Equal(t, tool, siblingNRJMXExec(executable))
	})

	t.Run("returns empty when there is no tool", func(t *testing.T) {
		executable := touch(t, resolvedTempDir(t), "integration")

		assert.Empty(t, siblingNRJMXExec(executable))
	})

	t.Run("ignores a directory named like the tool", func(t *testing.T) {
		dir := resolvedTempDir(t)
		executable := touch(t, dir, "integration")
		require.NoError(t, os.Mkdir(filepath.Join(dir, nrJMXFileName), 0o755))

		assert.Empty(t, siblingNRJMXExec(executable))
	})
}

func TestGetNRJMXExec(t *testing.T) {
	t.Run("environment variable takes precedence", func(t *testing.T) {
		t.Setenv(nrJMXEnvVar, "/custom/nrjmx")

		assert.Equal(t, "/custom/nrjmx", getNRJMXExec())
	})

	t.Run("falls back to the OS default", func(t *testing.T) {
		t.Setenv(nrJMXEnvVar, "")

		// The test binary has no nrjmx next to it.
		assert.Equal(t, defaultNRJMXExec, getNRJMXExec())
	})
}

// resolvedTempDir returns a temporary directory free of symlinks, as some OSes place them under a symlinked path.
func resolvedTempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	return dir
}

func touch(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, nil, 0o755))

	return path
}
