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
		if len(args) > 0 {
			message = fmt.Sprintf(message, args...)
		}
		return &nrprotocol.JMXConnectionError{
			Message: message,
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
func (p *Process) Start() (*Process, error) {
	if p.state.IsRunning() {
		return p, ErrProcessAlreadyRunning
	}

	p.cmd = exec.CommandContext(p.ctx, filepath.Clean(getNRJMXExec()), nrJMXV2Flag)

	var err error

	defer func() {
		if err != nil {
			_ = p.Terminate()
		}
	}()

	p.Stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return nil, NewJMXConnectionError("failed to create stdout pipe to %q: %v", p.cmd.Path, err)
	}

	p.Stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return nil, NewJMXConnectionError("failed to create stdin pipe to %q: %v", p.cmd.Path, err)
	}

	stderrBuff := NewDefaultLimitedBuffer()
	p.cmd.Stderr = stderrBuff

	if err := p.cmd.Start(); err != nil {
		return p, NewJMXConnectionError("failed to Start %q: %v", p.cmd.Path, err)
	}
	p.state.Start()

	go func() {
		err := p.cmd.Wait()
		if err != nil {
			err = NewJMXConnectionError("nrjmx Process exited with Error: %v: stderr: %s",
				err,
				stderrBuff.String())
		}
		p.state.Stop(err)
	}()

	return p, nil
}

// Error checks whether nrjmx subprocess reported any error.
func (p *Process) Error() error {
	select {
	case err, open := <-p.ErrorC():
		if err == nil && !open {
			// When the process exited with success but we call this function we report an error
			// to signal that a new query cannot be performed.
			return ErrProcessNotRunning
		}
		return err
	default:
		if !p.state.IsRunning() {
			return ErrProcessNotRunning
		}
		return nil
	}
}

// ErrorC returns jmx Process error channel.
func (p *Process) ErrorC() <-chan error {
	return p.state.ErrorC()
}

// WaitExit should be called when nrjmx process is expected to wait.
// It will ensure that the exit status will be captured.
// In case process doesn't exit in the expected time will be terminated.
func (p *Process) WaitExit(timeout time.Duration) error {
	select {
	case err := <-p.state.ErrorC():
		return err
	case <-time.After(timeout):
		err := p.Terminate()
		return NewJMXConnectionError(
			"timeout exceeded while waiting for nrjmx Process to exit gracefully, attempting to Terminate the Process, error: %v",
			err,
		)
	}
}

// Terminate will stop the nrjmx process.
func (p *Process) Terminate() (err error) {
	if !p.state.IsRunning() {
		return
	}

	if stdoutErr := p.Stdout.Close(); stdoutErr != nil {
		err = NewJMXConnectionError("failed to detach stdout from %q: %v", p.cmd.Path, stdoutErr)
	}
	if stdinErr := p.Stdin.Close(); stdinErr != nil {
		err = NewJMXConnectionError("failed to detach stdin from %q: %v", p.cmd.Path, stdinErr)
	}

	p.cancel()

	return err
}

// GetPID returns nrjmx subprocess pid.
func (p *Process) GetPID() int {
	if p.cmd == nil || p.cmd.Process == nil {
		return -1
	}
	return p.cmd.Process.Pid
}

// GetOSProcessState returns the os.ProcessState for nrjmx process.
func (p *Process) GetOSProcessState() *os.ProcessState {
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
