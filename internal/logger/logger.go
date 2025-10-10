package logger

// Logger is a minimal logging interface used by this project.
type Logger interface {
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Debug(args ...any)
	Sync() error
}
