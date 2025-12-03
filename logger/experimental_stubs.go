//go:build !experimental

/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package logger

// Stubs for experimental features when not building with experimental tag

// IsUsingSlog always returns false when experimental features are disabled
func (l *Logger) IsUsingSlog() bool {
	return false
}

// IsUsingSlog always returns false when experimental features are disabled
func IsUsingSlog() bool {
	return false
}
