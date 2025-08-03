package simplevisor

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics holds all OpenTelemetry metrics for the supervisor
type Metrics struct {
	// General status metrics
	processesRunning metric.Int64UpDownCounter
	processesTotal   metric.Int64UpDownCounter
	
	// Critical monitoring metrics
	processRestartCount    metric.Int64UpDownCounter
	processStatusGauge     metric.Int64UpDownCounter
	
	// Event counters
	processStarted     metric.Int64Counter
	processStopped     metric.Int64Counter
	processRestarted   metric.Int64Counter
	processPanics      metric.Int64Counter
	restartLimitExceeded metric.Int64Counter
	
	// Performance metrics
	shutdownTimeouts metric.Int64Counter
}

// newMetrics creates and initializes all metrics
func newMetrics() (*Metrics, error) {
	meter := otel.Meter("simplevisor")
	
	var err error
	m := &Metrics{}
	
	// General status metrics
	m.processesRunning, err = meter.Int64UpDownCounter(
		"simplevisor_processes_running",
		metric.WithDescription("Number of currently running processes"),
	)
	if err != nil {
		return nil, err
	}
	
	m.processesTotal, err = meter.Int64UpDownCounter(
		"simplevisor_processes_total",
		metric.WithDescription("Total number of registered processes by status"),
	)
	if err != nil {
		return nil, err
	}
	
	// Critical monitoring metrics
	m.processRestartCount, err = meter.Int64UpDownCounter(
		"simplevisor_process_restart_count",
		metric.WithDescription("Current restart count for each process"),
	)
	if err != nil {
		return nil, err
	}
	
	m.processStatusGauge, err = meter.Int64UpDownCounter(
		"simplevisor_process_status",
		metric.WithDescription("Process status (1=running, 0=stopped, -1=restarting)"),
	)
	if err != nil {
		return nil, err
	}
	
	// Event counters
	m.processStarted, err = meter.Int64Counter(
		"simplevisor_process_started_total",
		metric.WithDescription("Total number of process starts"),
	)
	if err != nil {
		return nil, err
	}
	
	m.processStopped, err = meter.Int64Counter(
		"simplevisor_process_stopped_total",
		metric.WithDescription("Total number of process stops"),
	)
	if err != nil {
		return nil, err
	}
	
	m.processRestarted, err = meter.Int64Counter(
		"simplevisor_process_restarted_total",
		metric.WithDescription("Total number of process restarts"),
	)
	if err != nil {
		return nil, err
	}
	
	m.processPanics, err = meter.Int64Counter(
		"simplevisor_process_panics_total",
		metric.WithDescription("Total number of process panics"),
	)
	if err != nil {
		return nil, err
	}
	
	m.restartLimitExceeded, err = meter.Int64Counter(
		"simplevisor_restart_limit_exceeded_total",
		metric.WithDescription("Total number of processes that exceeded restart limits"),
	)
	if err != nil {
		return nil, err
	}
	
	// Performance metrics
	m.shutdownTimeouts, err = meter.Int64Counter(
		"simplevisor_shutdown_timeouts_total",
		metric.WithDescription("Total number of shutdown timeouts"),
	)
	if err != nil {
		return nil, err
	}
	
	return m, nil
}

// recordProcessStarted records when a process starts
func (m *Metrics) recordProcessStarted(name string, policy RestartPolicy) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
		attribute.String("restart_policy", policy.String()),
	}
	
	m.processStarted.Add(context.Background(), 1, metric.WithAttributes(attrs...))
	m.processesRunning.Add(context.Background(), 1, metric.WithAttributes(attrs...))
	m.updateProcessStatus(name, StatusRunning)
}

// recordProcessStopped records when a process stops
func (m *Metrics) recordProcessStopped(name string, reason string) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
		attribute.String("reason", reason), // error, success, manual, shutdown
	}
	
	m.processStopped.Add(context.Background(), 1, metric.WithAttributes(attrs...))
	m.processesRunning.Add(context.Background(), -1, 
		metric.WithAttributes(attribute.String("process_name", name)))
	m.updateProcessStatus(name, StatusStopped)
}

// recordProcessRestarted records when a process restarts
func (m *Metrics) recordProcessRestarted(name string, policy RestartPolicy, restartCount int) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
		attribute.String("restart_policy", policy.String()),
	}
	
	m.processRestarted.Add(context.Background(), 1, metric.WithAttributes(attrs...))
	m.updateRestartCount(name, restartCount)
	m.updateProcessStatus(name, StatusRestarting)
}

// recordProcessPanic records when a process panics
func (m *Metrics) recordProcessPanic(name string) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
	}
	
	m.processPanics.Add(context.Background(), 1, metric.WithAttributes(attrs...))
}

// recordRestartLimitExceeded records when a process exceeds restart limits
func (m *Metrics) recordRestartLimitExceeded(name string, maxRestarts int) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
		attribute.Int("max_restarts", maxRestarts),
	}
	
	m.restartLimitExceeded.Add(context.Background(), 1, metric.WithAttributes(attrs...))
}

// recordShutdownTimeout records when shutdown times out
func (m *Metrics) recordShutdownTimeout() {
	if m == nil {
		return
	}
	
	m.shutdownTimeouts.Add(context.Background(), 1)
}

// updateProcessStatus updates the process status gauge
func (m *Metrics) updateProcessStatus(name string, status ProcessStatus) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
		attribute.String("status", status.String()),
	}
	
	var value int64
	switch status {
	case StatusRunning:
		value = 1
	case StatusStopped:
		value = 0
	case StatusRestarting:
		value = -1
	}
	
	// Reset previous status for this process
	m.resetProcessStatusGauge(name)
	
	// Set new status
	m.processStatusGauge.Add(context.Background(), value, metric.WithAttributes(attrs...))
}

// updateRestartCount updates the restart count gauge for a process
func (m *Metrics) updateRestartCount(name string, count int) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("process_name", name),
	}
	
	// Reset previous count and set new one
	// Note: This is a simplified approach. In practice, you might want to track
	// previous values to calculate the delta properly
	m.processRestartCount.Add(context.Background(), int64(count), 
		metric.WithAttributes(attrs...))
}

// updateTotalProcesses updates the total processes gauge
func (m *Metrics) updateTotalProcesses(count int, status ProcessStatus) {
	if m == nil {
		return
	}
	
	attrs := []attribute.KeyValue{
		attribute.String("status", status.String()),
	}
	
	m.processesTotal.Add(context.Background(), int64(count), 
		metric.WithAttributes(attrs...))
}

// resetProcessStatusGauge resets the status gauge for a process
// This is a helper to handle gauge updates properly
func (m *Metrics) resetProcessStatusGauge(name string) {
	if m == nil {
		return
	}
	
	// Reset all possible status values for this process to 0
	statuses := []ProcessStatus{StatusRunning, StatusStopped, StatusRestarting}
	
	for _, status := range statuses {
		attrs := []attribute.KeyValue{
			attribute.String("process_name", name),
			attribute.String("status", status.String()),
		}
		
		var resetValue int64
		switch status {
		case StatusRunning:
			resetValue = -1
		case StatusStopped:
			resetValue = 0
		case StatusRestarting:
			resetValue = 1
		}
		
		m.processStatusGauge.Add(context.Background(), resetValue, 
			metric.WithAttributes(attrs...))
	}
}

// String methods for enums to provide readable metric labels
func (r RestartPolicy) String() string {
	switch r {
	case RestartNever:
		return "never"
	case RestartAlways:
		return "always"
	case RestartOnFailure:
		return "on_failure"
	default:
		return "unknown"
	}
}

func (s ProcessStatus) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusRunning:
		return "running"
	case StatusRestarting:
		return "restarting"
	default:
		return "unknown"
	}
}