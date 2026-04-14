package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	service string
}

func New(service string) *Logger {
	return &Logger{service: service}
}

func (l *Logger) log(level, msg string, args ...interface{}) {
	entry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     level,
		"service":   l.service,
		"message":   msg,
	}
	// Pair up args as key-value context
	for i := 0; i+1 < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		entry[key] = args[i+1]
	}
	b, _ := json.Marshal(entry)
	fmt.Fprintln(os.Stdout, string(b))
}

func (l *Logger) Info(msg string, args ...interface{})  { l.log("INFO", msg, args...) }
func (l *Logger) Error(msg string, args ...interface{}) { l.log("ERROR", msg, args...) }
func (l *Logger) Warn(msg string, args ...interface{})  { l.log("WARN", msg, args...) }
func (l *Logger) Debug(msg string, args ...interface{}) { l.log("DEBUG", msg, args...) }

func (l *Logger) With(args ...interface{}) *Logger {
	return l // simplified — could extend to carry fields
}
