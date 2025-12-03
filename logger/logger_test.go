package logger

import (
	"testing"
)

func TestBackwardsCompatibility(t *testing.T) {
	// Test that all existing functions still work
	Print("test print")
	Printf("test printf %s", "value")
	Info("test info")
	Infof("test infof %s", "value")
	Error("test error")
	Errorf("test errorf %s", "value")
	Debug("test debug")
	Debugf("test debugf %s", "value")
	Warn("test warn")
	Warnf("test warnf %s", "value")
}

func TestNewLogger(t *testing.T) {
	// Test creating new logger instances
	logger1 := New()
	if logger1.IsUsingSlog() {
		t.Error("New logger should not be using slog by default")
	}

	// Test that logger methods work
	logger1.Info("test logger instance method")
	logger1.Error("test error", "key", "value")
}

func TestLoggerMethods(t *testing.T) {
	logger := New()

	// Test all logger instance methods
	logger.Print("test print")
	logger.Printf("test printf %s", "value")
	logger.Info("test info")
	logger.Infof("test infof %s", "value")
	logger.Error("test error")
	logger.Errorf("test errorf %s", "value")
	logger.Debug("test debug")
	logger.Debugf("test debugf %s", "value")
	logger.Warn("test warn")
	logger.Warnf("test warnf %s", "value")
}

func Example() {
	// Backwards compatible usage - unchanged
	Info("Application started")
	Error("Something went wrong", "error", "connection failed")

	// Instance usage
	logger := New()
	logger.Info("Instance logging")
	logger.Error("Instance error", "user_id", 123)
}
