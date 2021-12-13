package nrjmx

import (
	"context"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx/internal/nrprotocol"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	nrJMXEnvVar = "NR_JMX_TOOL"
	nrJMXV2Flag = "-v2"
)

var (
	NewJMXConnectionError = func(message string, args ...interface{}) *nrprotocol.JMXConnectionError {
		return &nrprotocol.JMXConnectionError{
			Message: fmt.Sprintf(message, args),
		}
	}
	ErrProcessAlreadyRunning = NewJMXConnectionError("nrjmx subprocess is already running")
	ErrProcessNotRunning     = NewJMXConnectionError("nrjmx subprocess is not running")
)

// Process will handle the nrjmx subprocess.
type Process struct {
	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	state  *ProcessState
}

// NewProcess returns a new nrjmx Process.
func NewProcess(ctx context.Context) *Process {
	ctx, cancel := context.WithCancel(ctx)

	return &Process{
		ctx:    ctx,
		cancel: cancel,
		state:  NewProcessState(),
	}
}

// Start the nrjmx subprocess.
func (n *Process) Start() (*Process, error) {
	if n.state.IsRunning() {
		return n, ErrProcessAlreadyRunning
	}

	n.cmd = exec.CommandContext(n.ctx, filepath.Clean(getNRJMXExec()), nrJMXV2Flag)

	var err error

	defer func() {
		if err != nil {
			_ = n.Terminate()
		}
	}()

	n.Stdout, err = n.cmd.StdoutPipe()
	if err != nil {
		return nil, NewJMXConnectionError("failed to create stdout pipe to %q: %v", n.cmd.Path, err)
	}

	n.Stdin, err = n.cmd.StdinPipe()
	if err != nil {
		return nil, NewJMXConnectionError("failed to create stdin pipe to %q: %v", n.cmd.Path, err)
	}

	stderrBuff := NewDefaultLimitedBuffer()
	n.cmd.Stderr = stderrBuff

	if err := n.cmd.Start(); err != nil {
		return n, NewJMXConnectionError("failed to Start %q: %v", n.cmd.Path, err)
	}
	n.state.Start()

	go func() {
		err := n.cmd.Wait()
		if err != nil {
			err = NewJMXConnectionError("nrjmx Process exited with Error: %w: stderr: %s",
				err,
				stderrBuff.String())
		}
		n.state.Stop(err)
	}()

	return n, nil
}

// Error checks whether nrjmx subprocess reported any error.
func (n *Process) Error() error {
	select {
	case err, open := <-n.state.ErrorC():
		if err == nil && !open {
			// When the process exited with success but we call this function we report an error
			// to signal that a new query cannot be performed.
			return ErrProcessNotRunning
		}
		return err
	default:
		if !n.state.IsRunning() {
			return ErrProcessNotRunning
		}
		return nil
	}
}

// ErrorC returns jmx Process error channel.
func (n *Process) ErrorC() <-chan error {
	return n.state.ErrorC()
}

// WaitExit should be called when nrjmx process is expected to wait.
// It will ensure that the exit status will be captured.
// In case process doesn't exit in the expected time will be terminated.
func (n *Process) WaitExit(timeout time.Duration) error {
	select {
	case err := <-n.state.ErrorC():
		return err
	case <-time.After(timeout):
		err := n.Terminate()
		return NewJMXConnectionError(
			"timeout exceeded while waiting for nrjmx Process to exit gracefully, attempting to Terminate the Process, error: %w",
			err,
		)
	}
}

// Terminate will stop the nrjmx process.
func (n *Process) Terminate() (err error) {
	if !n.state.IsRunning() {
		return
	}

	if stdoutErr := n.Stdout.Close(); stdoutErr != nil {
		err = NewJMXConnectionError("failed to detach stdout from %q: %w", n.cmd.Path, stdoutErr)
	}
	if stdinErr := n.Stdin.Close(); stdinErr != nil {
		err = NewJMXConnectionError("failed to detach stdin from %q: %w", n.cmd.Path, stdinErr)
	}

	n.cancel()

	return err
}

// GetPID returns nrjmx subprocess pid.
func (n *Process) GetPID() int {
	if n.cmd == nil || n.cmd.Process == nil {
		return -1
	}
	return n.cmd.Process.Pid
}

// GetOSProcessState returns the os.ProcessState for nrjmx process.
func (n *Process) GetOSProcessState() *os.ProcessState {
	if n.cmd == nil {
		return nil
	}
	return n.cmd.ProcessState
}

// getNRJMXExec returns the path to nrjmx tool.
func getNRJMXExec() string {
	nrJMXExec := os.Getenv(nrJMXEnvVar)
	if nrJMXExec == "" {
		nrJMXExec = defaultNRJMXExec
	}
	return nrJMXExec
}
