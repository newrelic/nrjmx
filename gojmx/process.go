package gojmx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var bufferSize = 1024 * 1024
var defaultNrjmxExec = "/usr/local/bin/nrjmx"

func getNrjmxExec() string {
	if os.Getenv("NR_JMX_TOOL") != "" {
		return os.Getenv("NR_JMX_TOOL")
	}
	return defaultNrjmxExec
}

// var defaultNrjmxExec = "/home/cristi/workspace/cppc/java/nrjmx/bin/nrjmx"

type jmxProcess struct {
	sync.Mutex
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
	Stdout       io.ReadCloser
	Stdin        io.WriteCloser
	errCh        chan error
	stderrBuffer *stderrBuffer
}

type stderrBuffer struct {
	cap  int
	buff bytes.Buffer
}

func (s *stderrBuffer) Write(p []byte) (int, error) {
	if len(p) > s.cap {
		p = p[len(p)-s.cap:]
	}
	if len(p)+s.buff.Len() > s.cap {
		data := s.buff.String()
		data = data[s.cap-len(p):]
		s.buff.Reset()
		_, err := s.buff.Write([]byte(data))
		if err != nil {
			return 0, err
		}
	}
	return s.buff.Write(p)
}

func (s *stderrBuffer) WriteString(p string) (int, error) {
	return s.Write([]byte(p))
}

func (s *stderrBuffer) String() string {
	return s.buff.String()
}

func NewStderrBuffer(capacity int) *stderrBuffer {
	return &stderrBuffer{
		cap: capacity,
	}
}

func (j *jmxProcess) Start() (*jmxProcess, error) {
	if j.IsRunning() {
		return j, ErrAlreadyStarted
	}

	var err error
	j.Stdout, err = j.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe to %q: %v", j.cmd.Path, err)
	}

	j.Stdin, err = j.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe to %q: %v", j.cmd.Path, err)
	}

	j.stderrBuffer = NewStderrBuffer(bufferSize)
	j.cmd.Stderr = j.stderrBuffer

	if err := j.cmd.Start(); err != nil {
		return j, fmt.Errorf("failed to start %q: %v", j.cmd.Path, err)
	}
	j.SetIsRunning(true)

	go func() {
		// stderr we must read before wait, not with strings builder
		err := j.cmd.Wait()
		if err != nil {
			j.errCh <- fmt.Errorf("%s: %w", j.stderrBuffer.String(), err)
		}
		j.SetIsRunning(false)
	}()

	return j, nil
}

func NewJMXProcess(ctx context.Context) *jmxProcess {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, filepath.Clean(getNrjmxExec()), "-v2")

	//cmd := exec.CommandContext(ctx, "java", "-cp", "/Users/cciutea/workspace/nr/int/nrjmx/bin/*", "org.newrelic.nrjmx.Application", "-v2")

	return &jmxProcess{
		running: false,
		cmd:     cmd,
		ctx:     ctx,
		cancel:  cancel,
		errCh:   make(chan error, 1),
		}
}

func (p *jmxProcess) IsRunning() bool {
	p.Lock()
	defer p.Unlock()
	return p.running
}

func (j *jmxProcess) SetIsRunning(running bool) {
	j.Lock()
	defer j.Unlock()
	j.running = running
}

func (p *jmxProcess) Error() error {
	select {
	case err := <-p.errCh:
		return err
	default:
		p.Lock()
		defer p.Unlock()
		if !p.running {
			return ErrNotRunning
		}
		return nil
	}
}

func (p *jmxProcess) WaitExitError(timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return errors.New("timeout exceeded while waiting for jmx process error")
	case err := <-p.errCh:
		return err
	}
}

func (p *jmxProcess) stop() error {
	var errors error

	if err := p.Stdout.Close(); err != nil {
		errors = fmt.Errorf("failed to detach stdout from %q: %w", p.cmd.Path, err)
	}
	if err := p.Stdin.Close(); err != nil {
		errors = fmt.Errorf("failed to detach stdin from %q: %w", p.cmd.Path, err)
	}
	p.cancel()
	//err := p.cmd.Wait()
	//if err != nil {
	//	errors = fmt.Errorf("command failed %q: %w", p.cmd.Path, err)
	//}
	return errors
}
