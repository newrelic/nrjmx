package nrjmx

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var bufferSize = 4 * 1024 // initial 4KB per line.
var defaultNrjmxExec = "/usr/local/bin/nrjmx"

func getNrjmxExec() string {
	if os.Getenv("NR_JMX_TOOL") != "" {
		return os.Getenv("NR_JMX_TOOL")
	}
	return defaultNrjmxExec
}

// var defaultNrjmxExec = "/home/cristi/workspace/cppc/java/nrjmx/bin/nrjmx"

type jmxProcess struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	Stderr io.ReadCloser
}

func startJMXProcess(ctx context.Context) (*jmxProcess, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, filepath.Clean(getNrjmxExec()), "-v2")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe to %q: %v", cmd.Path, err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe to %q: %v", cmd.Path, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe to %q: %v", cmd.Path, err)
	}

	go func() {
		reader := bufio.NewReaderSize(stderr, bufferSize)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			line, err := reader.ReadString('\n')
			// API needs re to allow stderr full read before closing
			if err != nil {
				if _, isAlreadyClosed := err.(*os.PathError); !isAlreadyClosed && err != io.EOF {
					fmt.Fprintf(os.Stderr, "error while reading stderr: '%v'", err)
					continue
				}
				return
			}
			fmt.Fprint(os.Stderr, line)
		}
	}()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %q: %v", cmd.Path, err)
	}

	return &jmxProcess{
		Stdout: stdout,
		Stdin:  stdin,
		Stderr: stderr,
		cmd:    cmd,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (p *jmxProcess) stop(timeout time.Duration) error {
	var errors error

	if err := p.Stdout.Close(); err != nil {
		errors = fmt.Errorf("failed to detach stdout from %q: %w", p.cmd.Path, err)
	}
	if err := p.Stdin.Close(); err != nil {
		errors = fmt.Errorf("failed to detach stdin from %q: %w", p.cmd.Path, err)
	}
	if err := p.Stderr.Close(); err != nil {
		errors = fmt.Errorf("failed to detach stder from %q: %w", p.cmd.Path, err)
	}
	p.cancel()
	err := p.cmd.Wait()
	if err != nil {
		errors = fmt.Errorf("command failed %q: %w", p.cmd.Path, err)
	}
	return errors
}
