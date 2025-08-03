// Package simplevisor provides a simple, lightweight process supervisor for managing
// long-running goroutines with automatic restart capabilities and graceful shutdown.
//
// # Overview
//
// Simplevisor manages multiple named processes (goroutines) with configurable restart
// policies, panic recovery, and coordinated shutdown. It's designed for applications
// that need reliable background process management without complex dependencies.
//
// # Basic Usage
//
//	supervisor := simplevisor.New(5*time.Second, logger)
//
//	// Register a simple process
//	supervisor.Register("worker", func(ctx context.Context) error {
//		for {
//			select {
//			case <-ctx.Done():
//				return ctx.Err() // Graceful shutdown
//			case work := <-workChan:
//				processWork(work)
//			}
//		}
//	})
//
//	// Start all processes
//	supervisor.Run()
//
//	// Wait for shutdown signal and cleanup
//	supervisor.WaitOnShutdownSignal(func() {
//		// Optional cleanup callback
//		cleanupResources()
//	})
//
// # Restart Policies
//
// Configure automatic restart behavior for processes:
//
//	// Never restart
//	supervisor.Register("one-shot", handler)
//
//	// Always restart (even on successful completion)
//	supervisor.Register("persistent", handler,
//		simplevisor.WithRestart(simplevisor.RestartAlways, 5, 2*time.Second))
//
//	// Only restart on failures/panics
//	supervisor.Register("resilient", handler,
//		simplevisor.WithRestart(simplevisor.RestartOnFailure, 3, 1*time.Second))
//
// # Panic Recovery
//
// Handle panics in processes with custom recovery logic:
//
//	supervisor.Register("risky-process", riskyHandler,
//		simplevisor.WithRecover(func(recovered interface{}) {
//			log.Printf("Process panicked: %v", recovered)
//			// Send alert, record metrics, etc.
//		}))
//
// # Manual Process Control
//
// Control processes programmatically during runtime:
//
//	// Restart a specific process
//	err := supervisor.RestartProcess("worker")
//
//	// Stop a specific process
//	err := supervisor.StopProcess("worker")
//
//	// Check process status
//	status, err := supervisor.GetProcessStatus("worker")
//	switch status {
//	case simplevisor.StatusRunning:
//		// Process is active
//	case simplevisor.StatusStopped:
//		// Process has stopped
//	case simplevisor.StatusRestarting:
//		// Process is restarting
//	}
//
// # Graceful Shutdown
//
// Simplevisor provides coordinated shutdown with timeout protection:
//
//	// Automatic shutdown on OS signals
//	supervisor.WaitOnShutdownSignal(nil)
//
//	// Manual shutdown
//	supervisor.Shutdown()
//
// During shutdown:
// 1. All process contexts are cancelled
// 2. Processes should handle ctx.Done() and return gracefully
// 3. Supervisor waits for all processes to finish (with timeout)
// 4. Optional teardown callback is executed
//
// # Context-Based Cancellation
//
// All processes receive a context for cancellation detection:
//
//	func workerProcess(ctx context.Context) error {
//		ticker := time.NewTicker(1 * time.Second)
//		defer ticker.Stop()
//
//		for {
//			select {
//			case <-ctx.Done():
//				// Cleanup and exit gracefully
//				cleanup()
//				return ctx.Err()
//			case <-ticker.C:
//				// Do periodic work
//				doWork()
//			}
//		}
//	}
//
// # Thread Safety
//
// Simplevisor is thread-safe for:
// - Process registration (before Run() is called)
// - Manual process control (RestartProcess, StopProcess)
// - Status queries (GetProcessStatus, IsRunning, ProcessCount)
//
// However, the supervisor itself should be used from the main goroutine,
// particularly for Run() and WaitOnShutdownSignal().
//
// # Best Practices
//
// 1. Register all processes before calling Run() (duplicate names will panic)
// 2. Use context cancellation for graceful shutdown in process handlers
// 3. Set appropriate restart limits to prevent infinite restart loops
// 4. Use panic recovery for critical processes that must stay running
// 5. Keep process handlers lightweight and delegate heavy work to other goroutines
// 6. Always handle ctx.Done() in process main loops
//
// # Configuration Options
//
// Restart Policy Options:
//   - RestartNever: Process runs once and stops (default)
//   - RestartAlways: Process restarts regardless of exit condition
//   - RestartOnFailure: Process restarts only on errors or panics
//
// Default Values:
//   - Shutdown timeout: 5 seconds
//   - Max restarts: 3
//   - Restart delay: 1 second
//
// # Error Handling
//
// Processes should return errors for failure conditions:
//
//	func databaseWorker(ctx context.Context) error {
//		db, err := connectDB()
//		if err != nil {
//			return fmt.Errorf("failed to connect: %w", err)
//		}
//		defer db.Close()
//
//		for {
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			default:
//				if err := processDBWork(db); err != nil {
//					return fmt.Errorf("db work failed: %w", err)
//				}
//			}
//		}
//	}
//
// Returning an error will trigger restart behavior based on the configured policy.
//
// # OpenTelemetry Metrics
//
// Simplevisor provides comprehensive OpenTelemetry metrics for monitoring:
//
//	supervisor := simplevisor.New(5*time.Second, logger)
//	
//	// Enable metrics (optional)
//	if err := supervisor.EnableMetrics(); err != nil {
//		log.Fatal(err)
//	}
//
// Key metrics include:
// - simplevisor_processes_running: Currently running processes
// - simplevisor_process_restart_count: Restart count per process  
// - simplevisor_process_status: Process status gauge (1=running, 0=stopped, -1=restarting)
// - simplevisor_process_started_total: Process start events
// - simplevisor_process_stopped_total: Process stop events (by reason)
// - simplevisor_process_panics_total: Process panic events
// - simplevisor_restart_limit_exceeded_total: Critical restart failures
//
// Metrics are automatically recorded when EnableMetrics() is called.
//
// # Logging
//
// Simplevisor uses structured logging (slog) and logs:
// - Process start/stop events
// - Restart attempts with counts and delays
// - Panic recovery details
// - Shutdown progress and completion
//
// Pass a custom logger to New() or use nil for default console output.
package simplevisor