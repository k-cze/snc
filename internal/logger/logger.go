package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the verbosity level
type LogLevel int

const (
	ERROR LogLevel = iota
	WARN
	INFO
	DEBUG
)

var (
	currentLevel LogLevel = INFO
	logger       *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

// SetLevel sets the logging level
func SetLevel(level LogLevel) {
	currentLevel = level
}

// SetLevelFromString sets the logging level from a string
func SetLevelFromString(level string) {
	switch level {
	case "error":
		SetLevel(ERROR)
	case "warn", "warning":
		SetLevel(WARN)
	case "info":
		SetLevel(INFO)
	case "debug":
		SetLevel(DEBUG)
	default:
		SetLevel(INFO)
	}
}

// formatMessage formats a log message with timestamp and level
func formatMessage(level string, component, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if component != "" {
		return fmt.Sprintf("[%s] %s [%s] %s", timestamp, level, component, message)
	}
	return fmt.Sprintf("[%s] %s %s", timestamp, level, message)
}

// Error logs an error message
func Error(component, message string, args ...interface{}) {
	if currentLevel >= ERROR {
		msg := fmt.Sprintf(message, args...)
		logger.Println(formatMessage("ERROR", component, msg))
	}
}

// Warn logs a warning message
func Warn(component, message string, args ...interface{}) {
	if currentLevel >= WARN {
		msg := fmt.Sprintf(message, args...)
		logger.Println(formatMessage("WARN", component, msg))
	}
}

// Info logs an info message
func Info(component, message string, args ...interface{}) {
	if currentLevel >= INFO {
		msg := fmt.Sprintf(message, args...)
		logger.Println(formatMessage("INFO", component, msg))
	}
}

// Debug logs a debug message
func Debug(component, message string, args ...interface{}) {
	if currentLevel >= DEBUG {
		msg := fmt.Sprintf(message, args...)
		logger.Println(formatMessage("DEBUG", component, msg))
	}
}

// Fatal logs a fatal error and exits
func Fatal(component, message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	logger.Println(formatMessage("FATAL", component, msg))
	os.Exit(1)
}

// Printf provides a simple printf-style logging for compatibility
func Printf(component, format string, args ...interface{}) {
	Info(component, format, args...)
}

// Println provides a simple println-style logging for compatibility
func Println(component, message string, args ...interface{}) {
	Info(component, message, args...)
}

// Progress logs progress information
func Progress(component, operation, item string, args ...interface{}) {
	if currentLevel >= INFO {
		msg := fmt.Sprintf(item, args...)
		logger.Printf("[%s] PROGRESS [%s] %s: %s\n", time.Now().Format("15:04:05"), component, operation, msg)
	}
}

// Success logs success information
func Success(component, message string, args ...interface{}) {
	if currentLevel >= INFO {
		msg := fmt.Sprintf(message, args...)
		logger.Println(formatMessage("SUCCESS", component, msg))
	}
}
