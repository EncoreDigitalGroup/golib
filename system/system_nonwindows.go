//go:build !windows

/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package system

import "os"

func RunningElevated() bool {
	return os.Geteuid() == 0
}
