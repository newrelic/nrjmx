/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/newrelic/nrjmx/gojmx/internal/nrjmx"
)

const (
	nrJMXEnvVar = "NR_JMX_TOOL"
	nrJMXV2Flag = "-v2"
)

var (
	errProcessAlreadyRunning = newJMXClientError("nrjmx subprocess is already running")
	errProcessNotRunning     = newJMXClientError("nrjmx subprocess is not running")

	// errPingTimeout returned if pingTimeout exceeded.
	errPingTimeout = newJMXConnectionError("could not establish communication with nrjmx: subprocess response timeout")
)

// process will handle the nrjmx subprocess.
type process struct {
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	state  *nrjmx.ProcessState
}

// newProcess returns a new nrjmx process.
func newProcess(ctx context.Context) *process {
	ctx, cancel := context.WithCancel(ctx)

	return &process{
		ctx:    ctx,
		cancel: cancel,
		state:  nrjmx.NewProcessState(),
	}
}

// start the nrjmx subprocess.
func (p *process) start() (proc *process, err error) {
	if p.state.IsRunning() {
		return p, errProcessAlreadyRunning
	}

	p.cmd = buildExecCommand(p.ctx)

	defer func() {
		if err != nil {
			_ = p.terminate()
		}
	}()

	p.Stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return nil, newJMXClientError("failed to create stdout pipe to %q: %v", p.cmd.Path, err)
	}

	p.Stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return nil, newJMXClientError("failed to create stdin pipe to %q: %v", p.cmd.Path, err)
	}

	stderrBuff := nrjmx.NewDefaultLimitedBuffer()
	p.cmd.Stderr = stderrBuff

	if err = p.cmd.Start(); err != nil {
		return p, newJMXClientError("failed to start %q: %v", p.cmd.Path, err)
	}
	p.state.Start()

	go func() {
		err := p.cmd.Wait()
		if err != nil {
			err = newJMXClientError("nrjmx process exited with error: %v: stderr: %s",
				err,
				stderrBuff.String())
		}

		p.terminate()
		p.state.Stop(err)
	}()

	return p, nil
}

// error checks whether nrjmx subprocess reported any error.
func (p *process) error() error {
	select {
	case err, open := <-p.state.ErrorC():
		if err == nil && !open {
			// When the process exited with success, but we call this function we report an error
			// to signal that a new query cannot be performed.
			return errProcessNotRunning
		}
		return err
	default:
		if !p.state.IsRunning() {
			return errProcessNotRunning
		}
		return nil
	}
}

// waitExit should be called when nrjmx process is expected to wait.
// It will ensure that the exit status will be captured.
// In case process doesn't exit in the expected time will be terminated.
func (p *process) waitExit(timeout time.Duration) error {
	select {
	case err := <-p.state.ErrorC():
		return err
	case <-time.After(timeout):
		err := p.terminate()
		return newJMXClientError(
			"timeout exceeded while waiting for nrjmx process to exit gracefully, attempting to terminate the process, error: %v",
			err,
		)
	}
}

// terminate will stop the nrjmx process.
func (p *process) terminate() (err error) {

	if p.Stdout != nil {
		if stdoutErr := p.Stdout.Close(); stdoutErr != nil {
			err = newJMXClientError("failed to detach stdout from %q: %v", p.cmd.Path, stdoutErr)
		}
	}

	if p.Stdin != nil {
		if stdinErr := p.Stdin.Close(); stdinErr != nil {
			err = newJMXClientError("failed to detach stdin from %q: %v", p.cmd.Path, stdinErr)
		}
	}

	p.cancel()

	return err
}

// getPID returns nrjmx subprocess pid.
func (p *process) getPID() int {
	if p.cmd == nil || p.cmd.Process == nil {
		return -1
	}
	return p.cmd.Process.Pid
}

// getOSProcessState returns the os.ProcessState for nrjmx process.
func (p *process) getOSProcessState() *os.ProcessState {
	if p.cmd == nil {
		return nil
	}
	return p.cmd.ProcessState
}

// getNRJMXExec returns the path to nrjmx tool.
func getNRJMXExec() string {
	nrJMXExec := os.Getenv(nrJMXEnvVar)
	if nrJMXExec == "" {
		nrJMXExec = defaultNRJMXExec
	}
	return nrJMXExec
}
