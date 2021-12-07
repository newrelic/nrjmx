package gojmx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
	cmd     *exec.Cmd
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	Stdout  io.ReadCloser
	Stdin   io.WriteCloser
	errCh   chan error
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

func startJMXProcess(ctx context.Context) (*jmxProcess, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, filepath.Clean(getNrjmxExec()), "-v2")

	//cmd := exec.CommandContext(ctx, "java", "-cp", "/Users/cciutea/workspace/nr/int/nrjmx/bin/*", "org.newrelic.nrjmx.Application", "-v2")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe to %q: %v", cmd.Path, err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe to %q: %v", cmd.Path, err)
	}

	stderrbuf := NewStderrBuffer(bufferSize)
	cmd.Stderr = stderrbuf

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %q: %v", cmd.Path, err)
	}
	errCh := make(chan error, 1)

	jmxProcess := &jmxProcess{
		Stdout:  stdout,
		Stdin:   stdin,
		running: true,
		cmd:     cmd,
		ctx:     ctx,
		cancel:  cancel,
		errCh:   errCh,
	}
	go func() {
		// stderr we must read before wait, not with strings builder
		err := cmd.Wait()
		if err != nil {
			errCh <- fmt.Errorf("%s: %w", stderrbuf.String(), err)
		}
		jmxProcess.Lock()
		defer jmxProcess.Unlock()
		jmxProcess.running = false
	}()

	return jmxProcess, nil
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
