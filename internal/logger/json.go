// Package logger provides logging implementations for the application.
package logger

import (
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// JSONLogger logs JSON with UTC timestamp and severity using zap.
type JSONLogger struct {
	logger *zap.Logger
}

var _ Logger = (*JSONLogger)(nil)

// NewJSONLogger creates a new JSON logger using zap.
func NewJSONLogger() (*JSONLogger, error) {
	encCfg := zapcore.EncoderConfig{
		TimeKey:    "ts",
		LevelKey:   "level",
		MessageKey: "msg",
		CallerKey:  "caller",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339Nano))
		},
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)
	logger := zap.New(core)
	return &JSONLogger{logger: logger}, nil
}

// Info logs informational messages.
func (l *JSONLogger) Info(args ...any) { l.logger.Sugar().Info(args...) }

// Warn logs warning messages.
func (l *JSONLogger) Warn(args ...any) { l.logger.Sugar().Warn(args...) }

// Error logs error messages.
func (l *JSONLogger) Error(args ...any) { l.logger.Sugar().Error(args...) }

// Debug logs debug messages.
func (l *JSONLogger) Debug(args ...any) { l.logger.Sugar().Debug(args...) }

// MarshalJSON is provided so JSONLogger can be safely marshaled if needed.
func (l *JSONLogger) MarshalJSON() ([]byte, error) { return json.Marshal("json-logger") }
