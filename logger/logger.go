/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package logger

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"os"
)

type Logger struct {
	*log.Logger
}

var Default = New()

func New() *Logger {
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

	return &Logger{l}
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
