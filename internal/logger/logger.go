package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger abstrae la interfaz de logging
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// StdLogger implements Logger using standard log
type StdLogger struct {
	infoLog   *log.Logger
	warnLog   *log.Logger
	errorLog  *log.Logger
	fatalLog  *log.Logger
	debugLog  *log.Logger
	debugMode bool
}

// NewStdLogger creates a new standard logger
func NewStdLogger(debugMode bool) *StdLogger {
	return &StdLogger{
		infoLog:   log.New(os.Stdout, "✓ ", log.Lshortfile),
		warnLog:   log.New(os.Stdout, "⚠ ", log.Lshortfile),
		errorLog:  log.New(os.Stderr, "✗ ", log.Lshortfile),
		fatalLog:  log.New(os.Stderr, "💥 ", log.Lshortfile),
		debugLog:  log.New(os.Stdout, "🐛 ", log.Lshortfile),
		debugMode: debugMode,
	}
}

// Info registra un mensaje informativo
func (sl *StdLogger) Info(msg string, args ...interface{}) {
	sl.infoLog.Println(fmt.Sprintf(msg, args...))
}

// Warn registra una advertencia
func (sl *StdLogger) Warn(msg string, args ...interface{}) {
	sl.warnLog.Println(fmt.Sprintf(msg, args...))
}

// Error registra un error
func (sl *StdLogger) Error(msg string, args ...interface{}) {
	sl.errorLog.Println(fmt.Sprintf(msg, args...))
}

// Fatal registra un error fatal y sale
func (sl *StdLogger) Fatal(msg string, args ...interface{}) {
	sl.fatalLog.Fatalln(fmt.Sprintf(msg, args...))
}

// Debug registra un mensaje de debug (solo si debugMode=true)
func (sl *StdLogger) Debug(msg string, args ...interface{}) {
	if sl.debugMode {
		sl.debugLog.Println(fmt.Sprintf(msg, args...))
	}
}

// Global logger instance
var globalLogger Logger = NewStdLogger(false)

// GetLogger devuelve la instancia global del logger
func GetLogger() Logger {
	return globalLogger
}

// SetLogger establece la instancia global del logger
func SetLogger(logger Logger) {
	globalLogger = logger
}
