/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type slogBackend struct {
	*slog.Logger
}

func (s *slogBackend) Print(msg string, keyvals ...interface{}) {
	s.Info(msg, keyvals...)
}

func (s *slogBackend) Printf(format string, keyvals ...interface{}) {
	s.Logger.Info(fmt.Sprintf(format, keyvals...))
}

func (s *slogBackend) Info(msg string, keyvals ...interface{}) {
	s.Logger.Info(msg, keyvals...)
}

func (s *slogBackend) Infof(format string, args ...interface{}) {
	s.Logger.Info(fmt.Sprintf(format, args...))
}

func (s *slogBackend) Error(msg string, keyvals ...interface{}) {
	s.Logger.Error(msg, keyvals...)
}

func (s *slogBackend) Errorf(format string, args ...interface{}) {
	s.Logger.Error(fmt.Sprintf(format, args...))
}

func (s *slogBackend) Debug(msg string, keyvals ...interface{}) {
	s.Logger.Debug(msg, keyvals...)
}

func (s *slogBackend) Debugf(format string, args ...interface{}) {
	s.Logger.Debug(fmt.Sprintf(format, args...))
}

func (s *slogBackend) Warn(msg string, keyvals ...interface{}) {
	s.Logger.Warn(msg, keyvals...)
}

func (s *slogBackend) Warnf(format string, args ...interface{}) {
	s.Logger.Warn(fmt.Sprintf(format, args...))
}

// NewWithSlog creates a new Logger instance with the provided slog.Logger as the backend
func NewWithSlog(slogger *slog.Logger) *Logger {
	return &Logger{
		backend: &slogBackend{slogger},
		slogger: slogger,
	}
}

// NewSlogJSON creates a new Logger instance with a JSON slog handler
func NewSlogJSON(w io.Writer) *Logger {
	return NewWithSlog(slog.New(slog.NewJSONHandler(w, nil)))
}

// NewSlogText creates a new Logger instance with a text slog handler
func NewSlogText(w io.Writer) *Logger {
	return NewWithSlog(slog.New(slog.NewTextHandler(w, nil)))
}

// SetSlog configures the Logger instance to use the provided slog.Logger as its backend
func (l *Logger) SetSlog(slogger *slog.Logger) {
	l.backend = &slogBackend{slogger}
	l.slogger = slogger
}

// GetSlog returns the underlying slog.Logger if one is configured, nil otherwise
func (l *Logger) GetSlog() *slog.Logger {
	if sl, ok := l.slogger.(*slog.Logger); ok {
		return sl
	}
	return nil
}

// IsUsingSlog returns true if the Logger is currently using an slog backend
func (l *Logger) IsUsingSlog() bool {
	return l.slogger != nil
}

// SetSlog configures the default Logger to use the provided slog.Logger as its backend
func SetSlog(slogger *slog.Logger) {
	Default.SetSlog(slogger)
}

// GetSlog returns the underlying slog.Logger from the default Logger if one is configured, nil otherwise
func GetSlog() *slog.Logger {
	return Default.GetSlog()
}

// IsUsingSlog returns true if the default Logger is currently using an slog backend
func IsUsingSlog() bool {
	return Default.IsUsingSlog()
}

// WithContext returns an slog.Logger with the provided context if the default Logger is using slog, nil otherwise
func WithContext(ctx context.Context) *slog.Logger {
	if sl := Default.GetSlog(); sl != nil {
		return sl.With(slog.Any("context", ctx))
	}
	return nil
}
