package nrjmx

import "sync"

// ProcessState will be used to signal the Error from exec.Command.Wait()
type ProcessState struct {
	sync.Mutex
	ch      chan error
	running bool
}

// NewProcessState returns a new ProcessState instance.
func NewProcessState() *ProcessState {
	return &ProcessState{}
}

// Start will be called after exec.Command.Start()
func (s *ProcessState) Start() {
	s.Lock()
	defer s.Unlock()
	if !s.running {
		s.ch = make(chan error, 1)
		s.running = true
	}
}

// ErrorC give access to the ProcessState Error channel.
func (s *ProcessState) ErrorC() <-chan error {
	return s.ch
}

// Stop is used to signal the state of exec.Command.Wait().
// Should be called immediately after exec.Command.Wait() with the Error resulted from Wait().
func (s *ProcessState) Stop(err error) {
	if err != nil {
		s.ch <- err
	}
	s.close()
}

// IsRunning returns if the ProcessState is in a running phase.
func (s *ProcessState) IsRunning() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}

// close will end the ProcessState.
func (s *ProcessState) close() {
	s.Lock()
	defer s.Unlock()
	if s.running && s.ch != nil {
		close(s.ch)
	}
	s.running = false
}
