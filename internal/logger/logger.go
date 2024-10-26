package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	// DEBUG level
	DEBUG LogLevel = iota
	// INFO level
	INFO
	// WARNING level
	WARNING
	// ERROR level
	ERROR
)

var levelStrings = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

// Logger represents a custom logger
type Logger struct {
	level     LogLevel
	logFile   *os.File
	logger    *log.Logger
	directory string
}

// NewLogger creates a new Logger instance
func NewLogger(level LogLevel, directory string) (*Logger, error) {
	logger := &Logger{
		level:     level,
		directory: directory,
	}

	// Ensure the log directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	err := logger.rotate()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}
	return logger, nil
}

// rotate creates a new log file for the current day
func (l *Logger) rotate() error {
	if l.logFile != nil {
		l.logFile.Close()
	}

	now := time.Now()
	filename := filepath.Join(l.directory, fmt.Sprintf("%s.log", now.Format("2006-01-02")))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}

	l.logFile = file
	l.logger = log.New(file, "", 0)
	return nil
}

// log writes a log message with the given level
func (l *Logger) log(level LogLevel, message string) {
	if level < l.level {
		return
	}

	now := time.Now()
	if now.Day() != time.Now().Day() {
		err := l.rotate()
		if err != nil {
			fmt.Printf("Error rotating log file: %v\n", err)
			return
		}
	}

	_, file, line, _ := runtime.Caller(2)
	logMessage := fmt.Sprintf("[%s] [%s] [%s:%d] %s",
		levelStrings[level],
		now.Format("2006-01-02 15:04:05"),
		filepath.Base(file),
		line,
		message)

	l.logger.Println(logMessage)
	fmt.Printf("\n%s", logMessage)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

// Warning logs a warning message
func (l *Logger) Warning(message string) {
	l.log(WARNING, message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}
