/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package logger

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Backend interface {
	Print(msg string, keyvals ...interface{})
	Printf(format string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Infof(format string, args ...interface{})
	Error(msg string, keyvals ...interface{})
	Errorf(format string, args ...interface{})
	Debug(msg string, keyvals ...interface{})
	Debugf(format string, args ...interface{})
	Warn(msg string, keyvals ...interface{})
	Warnf(format string, args ...interface{})
}

type Logger struct {
	backend Backend
	slogger interface{} // Stores *slog.Logger when using slog backend
}

type charmLogger struct {
	*log.Logger
}

func (c *charmLogger) Print(msg string, keyvals ...interface{}) {
	c.Logger.Print(msg, keyvals...)
}

func (c *charmLogger) Printf(format string, keyvals ...interface{}) {
	c.Logger.Printf(format, keyvals...)
}

func (c *charmLogger) Info(msg string, keyvals ...interface{}) {
	c.Logger.Info(msg, keyvals...)
}

func (c *charmLogger) Infof(format string, args ...interface{}) {
	c.Logger.Infof(format, args...)
}

func (c *charmLogger) Error(msg string, keyvals ...interface{}) {
	c.Logger.Error(msg, keyvals...)
}

func (c *charmLogger) Errorf(format string, args ...interface{}) {
	c.Logger.Errorf(format, args...)
}

func (c *charmLogger) Debug(msg string, keyvals ...interface{}) {
	c.Logger.Debug(msg, keyvals...)
}

func (c *charmLogger) Debugf(format string, args ...interface{}) {
	c.Logger.Debugf(format, args...)
}

func (c *charmLogger) Warn(msg string, keyvals ...interface{}) {
	c.Logger.Warn(msg, keyvals...)
}

func (c *charmLogger) Warnf(format string, args ...interface{}) {
	c.Logger.Warnf(format, args...)
}

var Default = New()

func New() *Logger {
	return newCharmLogger()
}

func newCharmLogger() *Logger {
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

	return &Logger{
		backend: &charmLogger{l},
		slogger: nil,
	}
}

func (l *Logger) Print(msg string, keyvals ...interface{}) {
	l.backend.Print(msg, keyvals...)
}

func (l *Logger) Printf(format string, keyvals ...interface{}) {
	l.backend.Printf(format, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.backend.Info(msg, keyvals...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.backend.Infof(format, args...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.backend.Error(msg, keyvals...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.backend.Errorf(format, args...)
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.backend.Debug(msg, keyvals...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.backend.Debugf(format, args...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.backend.Warn(msg, keyvals...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.backend.Warnf(format, args...)
}

func Print(msg string, keyvals ...interface{}) {
	Default.Print(msg, keyvals...)
}

func Printf(format string, keyvals ...interface{}) {
	Default.Printf(format, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	Default.Info(msg, keyvals...)
}

func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

func Error(msg string, keyvals ...interface{}) {
	Default.Error(msg, keyvals...)
}

func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}

func Debug(msg string, keyvals ...interface{}) {
	Default.Debug(msg, keyvals...)
}

func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

func Warn(msg string, keyvals ...interface{}) {
	Default.Warn(msg, keyvals...)
}

func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}
