/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// dualBackend writes to both charm logger (terminal) and slog (structured output) simultaneously
type dualBackend struct {
	charm *charmLogger
	slog  *slogBackend
}

func (d *dualBackend) Print(msg string, keyvals ...interface{}) {
	d.charm.Print(msg, keyvals...)
	d.slog.Print(msg, keyvals...)
}

func (d *dualBackend) Printf(format string, keyvals ...interface{}) {
	d.charm.Printf(format, keyvals...)
	d.slog.Printf(format, keyvals...)
}

func (d *dualBackend) Info(msg string, keyvals ...interface{}) {
	d.charm.Info(msg, keyvals...)
	d.slog.Info(msg, keyvals...)
}

func (d *dualBackend) Infof(format string, args ...interface{}) {
	d.charm.Infof(format, args...)
	d.slog.Infof(format, args...)
}

func (d *dualBackend) Error(msg string, keyvals ...interface{}) {
	d.charm.Error(msg, keyvals...)
	d.slog.Error(msg, keyvals...)
}

func (d *dualBackend) Errorf(format string, args ...interface{}) {
	d.charm.Errorf(format, args...)
	d.slog.Errorf(format, args...)
}

func (d *dualBackend) Debug(msg string, keyvals ...interface{}) {
	d.charm.Debug(msg, keyvals...)
	d.slog.Debug(msg, keyvals...)
}

func (d *dualBackend) Debugf(format string, args ...interface{}) {
	d.charm.Debugf(format, args...)
	d.slog.Debugf(format, args...)
}

func (d *dualBackend) Warn(msg string, keyvals ...interface{}) {
	d.charm.Warn(msg, keyvals...)
	d.slog.Warn(msg, keyvals...)
}

func (d *dualBackend) Warnf(format string, args ...interface{}) {
	d.charm.Warnf(format, args...)
	d.slog.Warnf(format, args...)
}

// NewDualLogger creates a logger that outputs to both charm (terminal) and slog (structured) backends
func NewDualLogger(slogger *slog.Logger) *Logger {
	charmBackend := &charmLogger{createStyledLogger()}
	slogBackend := &slogBackend{slogger}

	return &Logger{
		backend: &dualBackend{
			charm: charmBackend,
			slog:  slogBackend,
		},
		slogger: slogger,
	}
}

// NewDualLoggerJSON creates a dual logger with JSON structured output to the specified writer
func NewDualLoggerJSON(w io.Writer) *Logger {
	return NewDualLogger(slog.New(slog.NewJSONHandler(w, nil)))
}

// NewDualLoggerText creates a dual logger with text structured output to the specified writer
func NewDualLoggerText(w io.Writer) *Logger {
	return NewDualLogger(slog.New(slog.NewTextHandler(w, nil)))
}

// EnableDualOutput configures the Logger to output to both charm and slog backends
func (l *Logger) EnableDualOutput(slogger *slog.Logger) {
	charmBackend := &charmLogger{createStyledLogger()}
	slogBackend := &slogBackend{slogger}

	l.backend = &dualBackend{
		charm: charmBackend,
		slog:  slogBackend,
	}
	l.slogger = slogger
}

// EnableDualOutputJSON enables dual output with JSON structured logging to the specified writer
func (l *Logger) EnableDualOutputJSON(w io.Writer) {
	l.EnableDualOutput(slog.New(slog.NewJSONHandler(w, nil)))
}

// EnableDualOutputText enables dual output with text structured logging to the specified writer
func (l *Logger) EnableDualOutputText(w io.Writer) {
	l.EnableDualOutput(slog.New(slog.NewTextHandler(w, nil)))
}

// IsDualOutput returns true if the logger is configured for dual output
func (l *Logger) IsDualOutput() bool {
	_, isDual := l.backend.(*dualBackend)
	return isDual
}

// EnableDualOutputJSON enables dual output on the default logger with JSON structured logging
func EnableDualOutputJSON(w io.Writer) {
	Default.EnableDualOutputJSON(w)
}

// EnableDualOutputText enables dual output on the default logger with text structured logging
func EnableDualOutputText(w io.Writer) {
	Default.EnableDualOutputText(w)
}

// IsDualOutput returns true if the default logger is configured for dual output
func IsDualOutput() bool {
	return Default.IsDualOutput()
}

// createStyledLogger creates a charm logger with the default styling
func createStyledLogger() *log.Logger {
	styles := log.DefaultStyles()

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))

	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	l := log.New(os.Stdout)
	l.SetStyles(styles)
	return l
}

// ConfigureForCLI sets up logging for CLI applications:
// - Charm logger for user-facing terminal output
// - JSON structured logging to file for debugging
func ConfigureForCLI(logFilePath string) error {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	EnableDualOutputJSON(logFile)
	return nil
}

// ConfigureForMicroservice sets up logging based on environment:
// - In development: dual output (charm + structured file)
// - In production: structured logging only
func ConfigureForMicroservice(isDevelopment bool, logOutput io.Writer) {
	if isDevelopment {
		EnableDualOutputJSON(logOutput)
	} else {
		// Production: structured logging only
		slogger := slog.New(slog.NewJSONHandler(logOutput, nil))
		SetSlog(slogger)
	}
}

// ConfigureForEnvironment automatically configures logging based on common environment patterns
func ConfigureForEnvironment() {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "development", "dev":
		// Development: dual output to stderr for structured logs
		EnableDualOutputJSON(os.Stderr)
	case "production", "prod":
		// Production: structured JSON only
		slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		SetSlog(slogger)
	case "testing", "test":
		// Testing: charm only (no file output)
		Default = New()
	default:
		// Default: charm only for interactive use
		Default = New()
	}
}
