package main

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bufio"
	"context"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"thr/jmx"

	"github.com/apache/thrift/lib/go/thrift"
)

var bufferSize = 4 * 1024 // initial 4KB per line.
// var defaultNrjmxExec = "/usr/bin/nrjmx"
var defaultNrjmxExec = "/Users/cciutea/workspace/nr/demo_jmx/src/cppc/java/nrjmx/bin/nrjmx"

type JMXProcess struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	Stderr io.ReadCloser
}

func StartJMXProcess(ctx context.Context) (*JMXProcess, error) {
	ctx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(ctx, filepath.Clean(defaultNrjmxExec))

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

	return &JMXProcess{
		Stdout: stdout,
		Stdin:  stdin,
		Stderr: stderr,
		cmd:    cmd,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (p *JMXProcess) Stop() error {
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
		// panic(err)
	}
	return errors

}

func NewJMXServiceClient(jmxProcess *JMXProcess) (client jmx.JMXService, err error) {
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTJSONProtocolFactory()

	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTTransportFactory()

	var transport thrift.TTransport
	transport = thrift.NewStreamTransport(jmxProcess.Stdout, jmxProcess.Stdin)
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		return nil, err
	}

	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client = jmx.NewJMXServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return
}
