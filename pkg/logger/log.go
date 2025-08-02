// Package logger provides structured logging functionality
// built on top of Uber's zap logger.
//
// Example usage:
//
//	logger, err := logger.New(zapcore.InfoLevel)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer logger.Close()
package logger

import (
	"errors"
	"fmt"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

const (
	// logLayout is the layout of log time
	logLayout = "2006-01-02 15:04:05.000"
)

// NewString creates logger from string level (info, debug, error, etc.)
func NewString(level string) (*Logger, error) {
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}
	return New(l)
}

// NewDefault creates logger with INFO level
func NewDefault() *Logger {
	logger, _ := New(zapcore.InfoLevel)
	return logger
}

// New returns new instance of Logger
func New(l zapcore.Level) (*Logger, error) {
	lvl := zap.NewAtomicLevelAt(l)

	zl, err := zap.Config{
		Level:             lvl,
		Development:       false,
		Encoding:          "json",
		DisableStacktrace: true,
		DisableCaller:     true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			EncodeTime:     zapcore.TimeEncoderOfLayout(logLayout),
			EncodeDuration: zapcore.StringDurationEncoder,

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			NameKey:     "key",
			FunctionKey: zapcore.OmitKey,

			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		},
	}.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Logger{Logger: zl}, nil
}

// Close closes the logger
func (l *Logger) Close() error {
	if err := l.Logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EBADF) {
		return err
	}

	return nil
}
