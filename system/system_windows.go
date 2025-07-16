//go:build windows

/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package system

import "golang.org/x/sys/windows"

func RunningElevated() bool {
	var sid *windows.SID

	sid, _ = windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	token := windows.Token(0)
	isMember, _ := token.IsMember(sid)

	return isMember
}
