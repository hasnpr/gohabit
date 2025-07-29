package log

import (
	"errors"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger is the (logger) of zap
	Logger *zap.Logger
	// Level is the logger level of zap
	Level zap.AtomicLevel
)

// logLayout is the layout of log time
const logLayout = "2006-01-02 15:04:05.000"

// init: Set NewProduction as default logger. Config depend on Logger instance
func init() {
	var err error
	Level = zap.NewAtomicLevel()
	Logger, err = zap.Config{
		Level:             Level,
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
		panic(err)
	}
}

// CloseLogger closes the logger
func CloseLogger() error {
	if err := Logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EBADF) {
		return err
	}
	return nil
}
