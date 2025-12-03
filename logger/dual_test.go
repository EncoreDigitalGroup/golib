//go:build experimental

package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDualBackend(t *testing.T) {
	var structuredBuf bytes.Buffer
	logger := NewDualLoggerJSON(&structuredBuf)

	// Test that dual backend is active
	if !logger.IsDualOutput() {
		t.Error("Expected IsDualOutput() to return true for dual logger")
	}

	// Test that slog is also active
	if !logger.IsUsingSlog() {
		t.Error("Expected IsUsingSlog() to return true for dual logger")
	}

	// Test logging - should output to both backends
	logger.Info("test dual message", "key", "value")

	// Verify structured output was written
	if structuredBuf.Len() == 0 {
		t.Error("Expected structured output to be written")
	}

	// Verify it's JSON format
	output := structuredBuf.String()
	if !strings.Contains(output, `"msg":"test dual message"`) {
		t.Error("Expected JSON formatted output in structured buffer")
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Error("Expected key-value pairs in structured output")
	}
}

func TestDualLoggerMethods(t *testing.T) {
	var structuredBuf bytes.Buffer
	logger := NewDualLoggerText(&structuredBuf)

	// Test all logging methods
	logger.Print("print message")
	logger.Printf("printf message %s", "formatted")
	logger.Info("info message")
	logger.Infof("infof message %d", 42)
	logger.Error("error message")
	logger.Errorf("errorf message %s", "formatted")
	logger.Warn("warn message")
	logger.Warnf("warnf message %s", "formatted")
	logger.Debug("debug message")
	logger.Debugf("debugf message %s", "formatted")

	// Verify structured output was written
	if structuredBuf.Len() == 0 {
		t.Error("Expected structured output to be written")
	}

	output := structuredBuf.String()
	expectedMessages := []string{
		"print message",
		"printf message formatted",
		"info message",
		"infof message 42",
		"error message",
		"errorf message formatted",
		"warn message",
		"warnf message formatted",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected structured output to contain: %s", msg)
		}
	}
}

func TestEnableDualOutput(t *testing.T) {
	// Start with regular logger
	logger := New()
	if logger.IsDualOutput() {
		t.Error("New logger should not be dual output initially")
	}

	// Enable dual output
	var structuredBuf bytes.Buffer
	logger.EnableDualOutputJSON(&structuredBuf)

	// Verify dual output is now enabled
	if !logger.IsDualOutput() {
		t.Error("Expected IsDualOutput() to return true after EnableDualOutputJSON")
	}

	// Test logging
	logger.Info("test after dual enable")
	if structuredBuf.Len() == 0 {
		t.Error("Expected structured output after enabling dual output")
	}
}

func TestDefaultDualOutput(t *testing.T) {
	// Store original default to restore later
	original := Default

	// Test enabling dual output on default logger
	var structuredBuf bytes.Buffer
	EnableDualOutputJSON(&structuredBuf)

	if !IsDualOutput() {
		t.Error("Expected default logger to have dual output enabled")
	}

	// Test logging via package functions
	Info("test default dual", "key", "value")
	if structuredBuf.Len() == 0 {
		t.Error("Expected structured output from default dual logger")
	}

	// Restore original default
	Default = original
}

func TestConfigureForMicroservice(t *testing.T) {
	// Store original default to restore later
	original := Default

	// Test development configuration
	var devBuf bytes.Buffer
	ConfigureForMicroservice(true, &devBuf)

	if !IsDualOutput() {
		t.Error("Expected dual output in development mode")
	}

	Info("development message")
	if devBuf.Len() == 0 {
		t.Error("Expected structured output in development mode")
	}

	// Test production configuration
	var prodBuf bytes.Buffer
	ConfigureForMicroservice(false, &prodBuf)

	if IsDualOutput() {
		t.Error("Expected no dual output in production mode")
	}
	if !IsUsingSlog() {
		t.Error("Expected slog to be active in production mode")
	}

	Info("production message")
	if prodBuf.Len() == 0 {
		t.Error("Expected structured output in production mode")
	}

	// Restore original default
	Default = original
}

func TestConfigureForEnvironment(t *testing.T) {
	// Store original default and environment
	original := Default
	originalEnv := os.Getenv("ENVIRONMENT")

	// Test development environment
	os.Setenv("ENVIRONMENT", "development")
	ConfigureForEnvironment()
	if !IsDualOutput() {
		t.Error("Expected dual output for development environment")
	}

	// Test production environment
	os.Setenv("ENVIRONMENT", "production")
	ConfigureForEnvironment()
	if IsDualOutput() {
		t.Error("Expected no dual output for production environment")
	}
	if !IsUsingSlog() {
		t.Error("Expected slog for production environment")
	}

	// Test unknown environment (should default to charm only)
	os.Setenv("ENVIRONMENT", "unknown")
	ConfigureForEnvironment()
	if IsDualOutput() {
		t.Error("Expected no dual output for unknown environment")
	}
	if IsUsingSlog() {
		t.Error("Expected no slog for unknown environment")
	}

	// Restore original state
	Default = original
	os.Setenv("ENVIRONMENT", originalEnv)
}

func ExampleNewDualLogger() {
	// Create a logger that outputs to both terminal (charm) and structured file
	logFile, _ := os.Create("app.log")
	defer logFile.Close()

	logger := NewDualLoggerJSON(logFile)

	// This will show styled output in terminal AND write JSON to file
	logger.Info("User login", "user_id", 123, "ip", "192.168.1.1")
	logger.Error("Database error", "error", "connection timeout")
}

func ExampleConfigureForCLI() {
	// Set up logging for a CLI application
	err := ConfigureForCLI("cli-debug.log")
	if err != nil {
		panic(err)
	}

	// Users see styled terminal output, debug info goes to file
	Info("Processing files...")
	Error("Failed to process file", "file", "data.csv", "error", "permission denied")
}

func ExampleConfigureForMicroservice() {
	// Configure based on environment
	isDev := os.Getenv("ENVIRONMENT") == "development"

	var logOutput *os.File
	if isDev {
		logOutput = os.Stderr // Development: structured logs to stderr
	} else {
		logOutput = os.Stdout // Production: structured logs to stdout
	}

	ConfigureForMicroservice(isDev, logOutput)

	// In dev: both terminal styling and structured logs
	// In prod: structured JSON only
	Info("Server starting", "port", 8080)
	Error("Database connection failed", "retries", 3)
}
