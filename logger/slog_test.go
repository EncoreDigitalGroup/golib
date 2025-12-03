package logger

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestSlogIntegration(t *testing.T) {
	// Test that slog can be enabled
	var buf bytes.Buffer
	slogger := slog.New(slog.NewJSONHandler(&buf, nil))

	// Test setting slog on default logger
	SetSlog(slogger)

	// Verify slog is active
	if !IsUsingSlog() {
		t.Error("Expected IsUsingSlog() to return true")
	}

	// Test logging with slog backend
	Info("test slog message")

	// Verify output was written to buffer
	if buf.Len() == 0 {
		t.Error("Expected output to be written to slog buffer")
	}

	// Reset default logger
	Default = New()
}

func TestNewSlogLoggers(t *testing.T) {
	var buf bytes.Buffer

	// Test JSON logger
	jsonLogger := NewSlogJSON(&buf)
	if !jsonLogger.IsUsingSlog() {
		t.Error("NewSlogJSON logger should be using slog")
	}

	jsonLogger.Info("test json message")
	if buf.Len() == 0 {
		t.Error("Expected JSON output to be written to buffer")
	}

	// Verify it's JSON format
	output := buf.String()
	if !strings.Contains(output, `"msg":"test json message"`) {
		t.Error("Expected JSON formatted output")
	}

	// Test text logger
	buf.Reset()
	textLogger := NewSlogText(&buf)
	if !textLogger.IsUsingSlog() {
		t.Error("NewSlogText logger should be using slog")
	}

	textLogger.Info("test text message")
	if buf.Len() == 0 {
		t.Error("Expected text output to be written to buffer")
	}

	// Verify it's text format
	output = buf.String()
	if !strings.Contains(output, "test text message") {
		t.Error("Expected text formatted output")
	}
}

func TestNewWithSlog(t *testing.T) {
	var buf bytes.Buffer
	slogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger := NewWithSlog(slogger)
	if !logger.IsUsingSlog() {
		t.Error("NewWithSlog logger should be using slog")
	}

	// Test that we can get the underlying slog logger
	if logger.GetSlog() != slogger {
		t.Error("GetSlog should return the original slog logger")
	}

	logger.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Expected debug output to be written to buffer")
	}
}

func TestSlogLoggerMethods(t *testing.T) {
	var buf bytes.Buffer
	// Create logger with debug level enabled
	slogger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	logger := NewWithSlog(slogger)

	// Test all logger methods with slog backend
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

	// Verify output was written
	if buf.Len() == 0 {
		t.Error("Expected output to be written to slog buffer")
	}

	// Count the number of JSON lines (should be 10)
	lines := strings.Count(buf.String(), "\n")
	if lines < 10 {
		t.Errorf("Expected at least 10 log lines, got %d", lines)
	}
}

func TestSetSlogOnInstance(t *testing.T) {
	logger := New()
	if logger.IsUsingSlog() {
		t.Error("New logger should not be using slog initially")
	}

	var buf bytes.Buffer
	slogger := slog.New(slog.NewJSONHandler(&buf, nil))

	logger.SetSlog(slogger)
	if !logger.IsUsingSlog() {
		t.Error("Logger should be using slog after SetSlog")
	}

	if logger.GetSlog() != slogger {
		t.Error("GetSlog should return the set slog logger")
	}

	logger.Info("test after slog set")
	if buf.Len() == 0 {
		t.Error("Expected output after setting slog")
	}
}

func TestWithContext(t *testing.T) {
	// Test with no slog configured
	if WithContext(nil) != nil {
		t.Error("WithContext should return nil when no slog is configured")
	}

	// Test with slog configured
	var buf bytes.Buffer
	slogger := slog.New(slog.NewJSONHandler(&buf, nil))
	SetSlog(slogger)

	contextLogger := WithContext(nil)
	if contextLogger == nil {
		t.Error("WithContext should return a logger when slog is configured")
	}

	// Reset default logger
	Default = New()
}

func ExampleNewSlogJSON() {
	// Create a JSON structured logger
	var logOutput bytes.Buffer
	logger := NewSlogJSON(&logOutput)

	logger.Info("User login", "user_id", 123, "ip", "192.168.1.1")
	logger.Error("Database connection failed", "error", "timeout", "duration", "30s")

	// Output will be JSON formatted
}

func ExampleSetSlog() {
	// Enable structured logging on the default logger
	logFile, _ := os.Create("app.log")
	defer logFile.Close()

	slogger := slog.New(slog.NewJSONHandler(logFile, nil))
	SetSlog(slogger)

	// All existing logger calls now use structured logging
	Info("Application started", "version", "1.0.0")
	Error("Failed to process request", "request_id", "abc123", "error", "invalid input")
}
