package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

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
			if logger.zap == nil {
				t.Fatal("New() returned logger with nil zap instance")
			}
			if logger.level.Level() != tt.level {
				t.Errorf("New() level = %v, want %v", logger.level.Level(), tt.level)
			}

			// Clean up
			_ = logger.Close()
		})
	}
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

func TestLogger_LoggingOutput(t *testing.T) {
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

	// Log a message
	logger.zap.Info("test message", zapcore.Field{Key: "key", Type: zapcore.StringType, String: "value"})

	// Close writer and read output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify JSON structure
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level INFO, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "test message" {
		t.Errorf("Expected msg 'test message', got %v", logEntry["msg"])
	}
	if _, exists := logEntry["ts"]; !exists {
		t.Error("Timestamp field 'ts' missing from log output")
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

	// Log at different levels
	logger.zap.Debug("debug message")  // Should be filtered out
	logger.zap.Info("info message")    // Should be filtered out
	logger.zap.Warn("warn message")    // Should appear
	logger.zap.Error("error message")  // Should appear

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

func TestLogLayout(t *testing.T) {
	if logLayout != "2006-01-02 15:04:05.000" {
		t.Errorf("logLayout = %q, want %q", logLayout, "2006-01-02 15:04:05.000")
	}
}