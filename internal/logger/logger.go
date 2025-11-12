package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logWithLevel("INFO", msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logWithLevel("ERROR", msg, args...)
}

func (l *Logger) Success(msg string, args ...interface{}) {
	l.logWithLevel("SUCCESS", msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.logWithLevel("WARNING", msg, args...)
}

func (l *Logger) logWithLevel(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	prefix := l.getPrefix(level)
	formatted := fmt.Sprintf(msg, args...)
	l.Printf("%s %s %s", timestamp, prefix, formatted)
}

func (l *Logger) getPrefix(level string) string {
	switch level {
	case "INFO":
		return "ℹ️"
	case "ERROR":
		return "❌"
	case "SUCCESS":
		return "✅"
	case "WARNING":
		return "⚠️"
	default:
		return "📝"
	}
}