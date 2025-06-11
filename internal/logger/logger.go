package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	debug  bool
	logger *log.Logger
}

func New(debug bool) *Logger {
	return &Logger{
		debug:  debug,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		l.logger.Printf("[%s] [DEBUG] %s", timestamp, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("[%s] [INFO] %s", timestamp, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("[%s] [ERROR] %s", timestamp, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.logger.Printf("[%s] [FATAL] %s", timestamp, fmt.Sprintf(format, args...))
	os.Exit(1)
}
