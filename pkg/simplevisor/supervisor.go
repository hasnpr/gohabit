/*
Package simplevisor is a simple supervisor.
Supervisor Registers long-running processes. It runs long-running processes in go routine
and handles the panic with function registered by process or default recover.
Supervisor listen to shut-down signal and then runs all shutdown functions registered by processes.
*/
package simplevisor

import (
	"context"
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
	LogNSSupervisor                = "supervisor"
)

// ProcessFunc is a long-running process which listens on finishSignal.
type ProcessFunc func() error

// RecoverFunc is a function to execute when a process panics.
type RecoverFunc func(r any)

// ShutdownFunc is a function to execute for graceful shutdown.
type ShutdownFunc func(ctx context.Context)

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
	}
}

type Option func(p *Process)

type Process struct {
	name            string
	handler         ProcessFunc
	recoverHandler  RecoverFunc
	shutdownHandler ShutdownFunc
}

// WithRecover sets the recover handler for the process.
func WithRecover(handler RecoverFunc) Option {
	return func(p *Process) {
		p.recoverHandler = handler
	}
}

// WithShutdown sets the shutdown handler for the process.
func WithShutdown(handler ShutdownFunc) Option {
	return func(p *Process) {
		p.shutdownHandler = handler
	}
}

// Register registers a new process to supervisor.
// If the name isn't unique Register doesn't add process to the list.
func (s *Supervisor) Register(name string, handler ProcessFunc, options ...Option) {
	if s.checkIfNameAlreadyInUse(name) {
		s.logger.Warn("process name already in use, not registered", slog.String("process_name", name))
		return
	}

	process := Process{
		name:    name,
		handler: handler,
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
		go s.executeProcess(name, p)
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

func (s *Supervisor) executeProcess(name string, process Process) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("recover from panic", slog.String("process_name", name), slog.Any("panic", r))

			if process.recoverHandler != nil {
				process.recoverHandler(r)
			}
		}
	}()

	s.logger.Info("execute process", slog.String("process_name", name))

	err := process.handler()
	if err != nil {
		s.logger.Error("process execution finished", slog.String("process_name", name),
			slog.String("error", err.Error()))
	}
}

func (s *Supervisor) checkIfNameAlreadyInUse(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.processes[name]; ok {
		return true
	}

	return false
}

func (s *Supervisor) gracefulShutdown() {
	s.logger.Info("notify all processes to finish their jobs",
		slog.Duration("shutdown_timeout", s.shutdownTimeout),
		slog.Int("number_of_unfinished_processes", len(s.processes)))

	forceExitCtx, forceExitCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer forceExitCancel()

	s.lock.Lock()
	processes := make(map[string]Process, len(s.processes))
	maps.Copy(processes, s.processes)
	processCount := len(s.processes)
	s.lock.Unlock()

	// Use Waitgroup for proper coordination
	var wg sync.WaitGroup

	for n, p := range processes {
		wg.Add(1)
		go func(name string, process Process) {
			defer wg.Done()
			if process.shutdownHandler != nil {
				process.shutdownHandler(forceExitCtx)
			}

			s.logger.Info("process terminates gracefully", slog.String("process_name", name))
			s.removeProcess(name)
		}(n, p)
	}

	// Wait for all shutdowns to complete or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-forceExitCtx.Done():
	}

	s.logger.Info("supervisor terminates its job.",
		slog.Int("number_of_unfinished_processes", processCount))
}

func (s *Supervisor) removeProcess(name string) {
	s.lock.Lock()
	delete(s.processes, name)
	s.lock.Unlock()
}

func (s *Supervisor) IsRunning(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, exists := s.processes[name]

	return exists
}

func (s *Supervisor) ProcessCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.processes)
}
