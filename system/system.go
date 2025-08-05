/*
 * Copyright (c) 2025. Encore Digital Group.
 * All Rights Reserved.
 */

package system

import (
	"net"
	"runtime"
	"time"
)

const cloudflareDNS = "1.1.1.1:53"
const googleDNS = "8.8.8.8:53"

func InternetConnected() bool {
	dns := checkDNSConnection(cloudflareDNS)

	if !dns {
		dns = checkDNSConnection(googleDNS)
	}

	return dns
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsMac() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func checkDNSConnection(dns string) bool {
	conn, err := net.DialTimeout("tcp", dns, 5*time.Second)

	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}
