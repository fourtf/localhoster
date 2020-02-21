// +build windows

package main

import (
	"os"
)

const (
	hostsFilePath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

var (
	configFilePath = os.Getenv("USERPROFILE") + "\\localhoster.yaml"
)
