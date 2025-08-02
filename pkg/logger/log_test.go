package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		level zapcore.Level
	}{
		{"debug level", zapcore.DebugLevel},
		{"info level", zapcore.InfoLevel},
		{"warn level", zapcore.WarnLevel},
		{"error level", zapcore.ErrorLevel},
		{"fatal level", zapcore.FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}
			if logger == nil {
				t.Fatal("New() returned nil logger")
			}
			if logger.Logger == nil {
				t.Fatal("New() returned logger with nil embedded Logger")
			}

			// Test that embedded logger works
			logger.Info("test message")

			// Clean up
			_ = logger.Close()
		})
	}
}

func TestNewString(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		expectError bool
		expectedLvl zapcore.Level
	}{
		{"debug level", "debug", false, zapcore.DebugLevel},
		{"info level", "info", false, zapcore.InfoLevel},
		{"warn level", "warn", false, zapcore.WarnLevel},
		{"error level", "error", false, zapcore.ErrorLevel},
		{"fatal level", "fatal", false, zapcore.FatalLevel},
		{"panic level", "panic", false, zapcore.PanicLevel},
		{"invalid level", "invalid", true, zapcore.InfoLevel},
		{"uppercase level", "INFO", false, zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewString(tt.level)

			if tt.expectError {
				if err == nil {
					t.Errorf("NewString() expected error for level %q, got nil", tt.level)
				}
				if logger != nil {
					t.Errorf("NewString() expected nil logger for invalid level, got %v", logger)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewString() failed: %v", err)
			}
			if logger == nil {
				t.Fatal("NewString() returned nil logger")
			}
			if logger.Logger == nil {
				t.Fatal("NewString() returned logger with nil embedded Logger")
			}

			// Clean up
			_ = logger.Close()
		})
	}
}

func TestNewDefault(t *testing.T) {
	logger := NewDefault()
	if logger == nil {
		t.Fatal("NewDefault() returned nil logger")
	}
	if logger.Logger == nil {
		t.Fatal("NewDefault() returned logger with nil embedded Logger")
	}

	// Test that it works
	logger.Info("test message")

	// Clean up
	_ = logger.Close()
}

func TestLogger_Close(t *testing.T) {
	logger, err := New(zapcore.InfoLevel)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	err = logger.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}

func TestLogger_CloseMultipleTimes(t *testing.T) {
	logger, err := New(zapcore.InfoLevel)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// First close should work
	err = logger.Close()
	if err != nil {
		t.Errorf("First Close() failed: %v", err)
	}

	// Second close should also work (should not panic or error)
	err = logger.Close()
	if err != nil {
		t.Errorf("Second Close() failed: %v", err)
	}
}

func TestLogger_EmbeddedMethods(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	logger, err := New(zapcore.InfoLevel)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer logger.Close()

	// Test that embedded zap.Logger methods work directly
	logger.Info("info message", zap.String("key", "value"))
	logger.Error("error message", zap.Int("count", 42))

	// Close writer and read output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(lines))
	}

	// Check first log entry (info)
	var logEntry1 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &logEntry1); err != nil {
		t.Fatalf("Failed to parse first JSON output: %v", err)
	}
	if logEntry1["level"] != "INFO" {
		t.Errorf("Expected level INFO, got %v", logEntry1["level"])
	}
	if logEntry1["msg"] != "info message" {
		t.Errorf("Expected msg 'info message', got %v", logEntry1["msg"])
	}
	if logEntry1["key"] != "value" {
		t.Errorf("Expected key 'value', got %v", logEntry1["key"])
	}

	// Check second log entry (error)
	var logEntry2 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &logEntry2); err != nil {
		t.Fatalf("Failed to parse second JSON output: %v", err)
	}
	if logEntry2["level"] != "ERROR" {
		t.Errorf("Expected level ERROR, got %v", logEntry2["level"])
	}
	if logEntry2["msg"] != "error message" {
		t.Errorf("Expected msg 'error message', got %v", logEntry2["msg"])
	}
	if logEntry2["count"] != float64(42) { // JSON unmarshals numbers as float64
		t.Errorf("Expected count 42, got %v", logEntry2["count"])
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	// Create logger with WARN level
	logger, err := New(zapcore.WarnLevel)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer logger.Close()

	// Log at different levels using embedded methods
	logger.Debug("debug message") // Should be filtered out
	logger.Info("info message")   // Should be filtered out
	logger.Warn("warn message")   // Should appear
	logger.Error("error message") // Should appear

	// Close writer and read output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(lines))
	}

	// Check that only WARN and ERROR messages appear
	for _, line := range lines {
		if line == "" {
			continue
		}
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}
		level := logEntry["level"].(string)
		if level != "WARN" && level != "ERROR" {
			t.Errorf("Unexpected log level: %s", level)
		}
	}
}

func TestLogger_WithFields(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	logger, err := New(zapcore.InfoLevel)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer logger.Close()

	// Test With method from embedded logger
	childLogger := logger.With(zap.String("service", "test"), zap.Int("version", 1))
	childLogger.Info("child logger message")

	// Close writer and read output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check that fields are present
	if logEntry["service"] != "test" {
		t.Errorf("Expected service 'test', got %v", logEntry["service"])
	}
	if logEntry["version"] != float64(1) {
		t.Errorf("Expected version 1, got %v", logEntry["version"])
	}
}

func TestLogLayout(t *testing.T) {
	if logLayout != "2006-01-02 15:04:05.000" {
		t.Errorf("logLayout = %q, want %q", logLayout, "2006-01-02 15:04:05.000")
	}
}
