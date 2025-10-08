package logger

import (
	"encoding/json"
	"fmt"
	"os"
)

// PlainLogger writes human readable logs without timestamp or severity.
type PlainLogger struct{}

var _ Logger = PlainLogger{}

func (PlainLogger) Info(args ...any)  { fmt.Fprintln(os.Stdout, fmt.Sprint(args...)) }
func (PlainLogger) Warn(args ...any)  { fmt.Fprintln(os.Stdout, fmt.Sprint(args...)) }
func (PlainLogger) Error(args ...any) { fmt.Fprintln(os.Stderr, fmt.Sprint(args...)) }
func (PlainLogger) Debug(args ...any) { fmt.Fprintln(os.Stdout, fmt.Sprint(args...)) }

// MarshalJSON is provided so PlainLogger can be safely marshaled if needed.
func (PlainLogger) MarshalJSON() ([]byte, error) { return json.Marshal("plain-logger") }
