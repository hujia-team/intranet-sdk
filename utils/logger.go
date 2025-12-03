// Package utils provides utility functions for the MiniEye Intranet SDK.
package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the logging level.
type LogLevel int

// Define log levels
const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelNone:
		return "NONE"
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger for the MiniEye Intranet SDK.
type Logger struct {
	level     LogLevel
	errLog    *log.Logger
	outLog    *log.Logger
	component string
}

// DefaultLogger is the default logger instance.
var DefaultLogger *Logger

func init() {
	// Initialize default logger with INFO level
	DefaultLogger = NewLogger(LogLevelInfo)
}

// NewLogger creates a new logger with the specified log level.
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:     level,
		errLog:    log.New(os.Stderr, "", log.LstdFlags),
		outLog:    log.New(os.Stdout, "", log.LstdFlags),
		component: "default",
	}
}

// WithComponent sets the component name for the logger.
func (l *Logger) WithComponent(component string) *Logger {
	l.component = component
	return l
}

// SetLogLevel sets the log level for the logger.
func (l *Logger) SetLogLevel(level LogLevel) {
	l.level = level
}

// SetLogLevelFromString sets the log level from a string.
func (l *Logger) SetLogLevelFromString(level string) error {
	switch strings.ToUpper(level) {
	case "NONE":
		l.level = LogLevelNone
	case "ERROR":
		l.level = LogLevelError
	case "WARN":
		l.level = LogLevelWarn
	case "INFO":
		l.level = LogLevelInfo
	case "DEBUG":
		l.level = LogLevelDebug
	case "TRACE":
		l.level = LogLevelTrace
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}
	return nil
}

// Error logs an error message with component information.
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= LogLevelError {
		msg := fmt.Sprintf(format, args...)
		l.errLog.Printf("%s [ERROR] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), l.component, msg)
	}
}

// Warn logs a warning message with component information.
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level >= LogLevelWarn {
		msg := fmt.Sprintf(format, args...)
		l.outLog.Printf("%s [WARN] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), l.component, msg)
	}
}

// Info logs an info message with component information.
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		l.outLog.Printf("%s [INFO] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), l.component, msg)
	}
}

// Debug logs a debug message with component information.
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		l.outLog.Printf("%s [DEBUG] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), l.component, msg)
	}
}

// Trace logs a trace message with component information.
func (l *Logger) Trace(format string, args ...interface{}) {
	if l.level >= LogLevelTrace {
		msg := fmt.Sprintf(format, args...)
		l.outLog.Printf("%s [TRACE] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), l.component, msg)
	}
}

// LogRequest logs an HTTP request with details.
func (l *Logger) LogRequest(method, path string, statusCode int, latency time.Duration) {
	if l.level >= LogLevelInfo {
		level := "INFO"
		if statusCode >= 500 {
			level = "ERROR"
		} else if statusCode >= 400 {
			level = "WARN"
		}
		l.outLog.Printf("%s [%s] [%s] %s %s %d %v", 
			time.Now().Format("2006-01-02 15:04:05"), 
			level, 
			l.component, 
			method, 
			path, 
			statusCode, 
			latency)
	}
}

// Global log functions that use the default logger

// SetDefaultLogLevel sets the log level for the default logger.
func SetDefaultLogLevel(level LogLevel) {
	DefaultLogger.SetLogLevel(level)
}

// SetDefaultLogLevelFromString sets the log level for the default logger from a string.
func SetDefaultLogLevelFromString(level string) error {
	return DefaultLogger.SetLogLevelFromString(level)
}

// LogRequest logs an HTTP request with details using the default logger.
func LogRequest(method, path string, statusCode int, latency time.Duration) {
	DefaultLogger.LogRequest(method, path, statusCode, latency)
}

// Error logs an error message using the default logger.
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}

// Warn logs a warning message using the default logger.
func Warn(format string, args ...interface{}) {
	DefaultLogger.Warn(format, args...)
}

// Info logs an info message using the default logger.
func Info(format string, args ...interface{}) {
	DefaultLogger.Info(format, args...)
}

// Debug logs a debug message using the default logger.
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}

// Trace logs a trace message using the default logger.
func Trace(format string, args ...interface{}) {
	DefaultLogger.Trace(format, args...)
}
