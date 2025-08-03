/*
Package simplevisor is a simple supervisor.
Supervisor Registers long-running processes. It runs long-running processes in go routine
and handles the panic with function registered by process or default recover.
Supervisor listen to shut-down signal and then runs all shutdown functions registered by processes.
*/
package simplevisor

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	DefaultGracefulShutdownTimeout = 5 * time.Second
	DefaultRestartDelay            = 1 * time.Second
	DefaultMaxRestarts             = 3
	LogNSSupervisor                = "supervisor"
)

// RestartPolicy defines when a process should be restarted
type RestartPolicy int

const (
	RestartNever     RestartPolicy = iota // Never restart the process
	RestartAlways                         // Always restart the process
	RestartOnFailure                      // Only restart on error/panic
)

// ProcessStatus represents the current state of a process
type ProcessStatus int

const (
	StatusStopped ProcessStatus = iota
	StatusRunning
	StatusRestarting
)

// ProcessFunc is a long-running process which listens on context cancellation.
type ProcessFunc func(ctx context.Context) error

// RecoverFunc is a function to execute when a process panics.
type RecoverFunc func(r any)

// Supervisor is responsible to manage long-running processes.
// Supervisor is not for concurrent use and should be used as the main goroutine of app.
type Supervisor struct {
	shutDownCtx     context.Context
	shutDownCancel  context.CancelFunc
	logger          *slog.Logger
	lock            sync.Mutex
	processes       map[string]Process
	shutdownSignal  chan os.Signal
	shutdownTimeout time.Duration
	processWg       sync.WaitGroup  // Tracks running process goroutines
	metrics         metricsRecorder // Metrics recorder (NoOp by default, OpenTelemetry when enabled)
}

// New returns new instance of Supervisor.
func New(shutdownTimeout time.Duration, sLog *slog.Logger) *Supervisor {
	ctx, cancel := context.WithCancel(context.Background())

	if sLog == nil {
		sLog = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout), nil))
	}

	if shutdownTimeout <= 0 {
		shutdownTimeout = DefaultGracefulShutdownTimeout
	}

	return &Supervisor{
		shutDownCtx:     ctx,
		shutDownCancel:  cancel,
		lock:            sync.Mutex{},
		logger:          sLog.WithGroup(LogNSSupervisor),
		processes:       make(map[string]Process),
		shutdownSignal:  make(chan os.Signal, 1),
		shutdownTimeout: shutdownTimeout,
		metrics:         &noOpMetrics{}, // Default to NoOp metrics to avoid nil pointer issues
	}
}

// EnableMetrics initializes OpenTelemetry metrics for the supervisor.
// This is optional and should be called before registering processes for best results.
func (s *Supervisor) EnableMetrics() error {
	metrics, err := newMetrics()
	if err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}
	s.metrics = metrics
	return nil
}

type Option func(p *Process)

type Process struct {
	name           string
	handler        ProcessFunc
	recoverHandler RecoverFunc
	restartPolicy  RestartPolicy
	maxRestarts    int
	restartDelay   time.Duration
	restartCount   int
	status         ProcessStatus
}

// WithRecover sets the recover handler for the process.
func WithRecover(handler RecoverFunc) Option {
	return func(p *Process) {
		p.recoverHandler = handler
	}
}

// WithRestart sets the restart policy for the process.
func WithRestart(policy RestartPolicy, maxRestarts int, delay time.Duration) Option {
	return func(p *Process) {
		p.restartPolicy = policy
		p.maxRestarts = maxRestarts
		if delay <= 0 {
			delay = DefaultRestartDelay
		}
		p.restartDelay = delay
	}
}

// Register registers a new process to supervisor.
// Panics if the name isn't unique.
func (s *Supervisor) Register(name string, handler ProcessFunc, options ...Option) {
	s.panicIfNameAlreadyInUse(name)

	process := Process{
		name:         name,
		handler:      handler,
		maxRestarts:  DefaultMaxRestarts,
		restartDelay: DefaultRestartDelay,
		status:       StatusStopped,
	}

	for _, option := range options {
		option(&process)
	}

	s.lock.Lock()
	s.processes[name] = process
	s.lock.Unlock()
}

// Run spawns a new goroutine for each process.
// Spawned goroutine is responsible to handle the panic.
func (s *Supervisor) Run() {
	s.lock.Lock()
	processes := make(map[string]Process, len(s.processes))
	maps.Copy(processes, s.processes)
	s.lock.Unlock()

	// there is no need to use a goroutine pool such as Ants because this goroutine is long-running.
	for name, p := range processes {
		s.processWg.Add(1)
		go s.executeProcessWithRestart(name, p)
	}
}

func (s *Supervisor) Context() context.Context {
	return s.shutDownCtx
}

// WaitOnShutdownSignal wait to receive shutdown signal.
// WaitOnShutdownSignal should not be called in other goroutines except main goroutine of app.
// teardown is a callback function and will run at the last stage.
func (s *Supervisor) WaitOnShutdownSignal(teardown func()) {
	signal.Notify(s.shutdownSignal, os.Interrupt, syscall.SIGTERM)
	<-s.shutdownSignal

	s.shutdown(teardown)
}

// Shutdown manually shuts down the supervisor goroutine
func (s *Supervisor) Shutdown() {
	s.shutdown(nil)
}

func (s *Supervisor) shutdown(teardown func()) {
	s.gracefulShutdown()

	if teardown != nil {
		teardown()
	}

	s.shutDownCancel()
}

func (s *Supervisor) executeProcessWithRestart(name string, process Process) {
	defer s.processWg.Done()

	for {
		select {
		case <-s.shutDownCtx.Done():
			return
		default:
		}

		shouldRestart := s.executeProcess(name, process)

		if !shouldRestart {
			s.setProcessStatus(name, StatusStopped)
			return
		}

		// Increment restart count
		s.incrementRestartCount(name)

		// Check if we've exceeded max restarts
		if process.maxRestarts > 0 && s.getRestartCount(name) >= process.maxRestarts {
			s.logger.Error("process exceeded max restarts",
				slog.String("process_name", name),
				slog.Int("restart_count", s.getRestartCount(name)))
			s.metrics.recordRestartLimitExceeded(name, process.maxRestarts)
			s.setProcessStatus(name, StatusStopped)
			return
		}

		s.setProcessStatus(name, StatusRestarting)
		restartCount := s.getRestartCount(name)
		s.logger.Info("restarting process",
			slog.String("process_name", name),
			slog.Duration("delay", process.restartDelay),
			slog.Int("restart_count", restartCount))
		s.metrics.recordProcessRestarted(name, process.restartPolicy, restartCount)

		// Wait for restart delay or shutdown signal
		select {
		case <-time.After(process.restartDelay):
		case <-s.shutDownCtx.Done():
			return
		}
	}
}

func (s *Supervisor) executeProcess(name string, process Process) bool {
	var processErr error
	var panicOccurred bool

	defer func() {
		if r := recover(); r != nil {
			panicOccurred = true
			s.logger.Error("recover from panic", slog.String("process_name", name), slog.Any("panic", r))
			s.metrics.recordProcessPanic(name)

			if process.recoverHandler != nil {
				process.recoverHandler(r)
			}
		}
	}()

	s.logger.Info("execute process", slog.String("process_name", name))
	s.setProcessStatus(name, StatusRunning)
	s.metrics.recordProcessStarted(s.shutDownCtx, name, process.restartPolicy)

	processErr = process.handler(s.shutDownCtx)
	if processErr != nil {
		s.logger.Error("process execution finished", slog.String("process_name", name),
			slog.String("error", processErr.Error()))
		s.metrics.recordProcessStopped(name, "error")
	} else {
		s.metrics.recordProcessStopped(name, "success")
	}

	// Determine if we should restart based on policy
	switch process.restartPolicy {
	case RestartNever:
		return false
	case RestartAlways:
		return true
	case RestartOnFailure:
		return processErr != nil || panicOccurred
	default:
		return false
	}
}

func (s *Supervisor) panicIfNameAlreadyInUse(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.processes[name]; ok {
		s.logger.Error("process name already in use", slog.String("process_name", name))
		panic(fmt.Sprintf("process name %q already in use", name))
	}
}

func (s *Supervisor) gracefulShutdown() {
	s.logger.Info("notify all processes to finish their jobs",
		slog.Duration("shutdown_timeout", s.shutdownTimeout),
		slog.Int("number_of_processes", len(s.processes)))

	// Cancel context to signal all processes to shutdown
	s.shutDownCancel()

	// Wait for all process goroutines to finish with timeout
	done := make(chan struct{})
	go func() {
		s.processWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("all processes terminated gracefully")
	case <-time.After(s.shutdownTimeout):
		s.logger.Warn("shutdown timeout exceeded, some processes may still be running")
		s.metrics.recordShutdownTimeout()
	}

	s.logger.Info("supervisor terminates its job.")
}

func (s *Supervisor) IsRunning(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	process, exists := s.processes[name]
	if !exists {
		return false
	}

	return process.status == StatusRunning
}

func (s *Supervisor) ProcessCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.processes)
}

// GetProcessStatus returns the current status of a process
func (s *Supervisor) GetProcessStatus(name string) (ProcessStatus, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	process, exists := s.processes[name]
	if !exists {
		return StatusStopped, fmt.Errorf("process %s not found", name)
	}

	return process.status, nil
}

// setProcessStatus updates the status of a process
func (s *Supervisor) setProcessStatus(name string, status ProcessStatus) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if process, exists := s.processes[name]; exists {
		process.status = status
		s.processes[name] = process
	}
}

// incrementRestartCount increments the restart counter for a process
func (s *Supervisor) incrementRestartCount(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if process, exists := s.processes[name]; exists {
		process.restartCount++
		s.processes[name] = process
	}
}

// getRestartCount returns the current restart count for a process
func (s *Supervisor) getRestartCount(name string) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	if process, exists := s.processes[name]; exists {
		return process.restartCount
	}
	return 0
}

// resetRestartCount resets the restart counter for a process
func (s *Supervisor) resetRestartCount(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if process, exists := s.processes[name]; exists {
		process.restartCount = 0
		s.processes[name] = process
	}
}
