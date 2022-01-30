//go:build (linux && !android && !nox11) || freebsd || openbsd || !windows
// +build linux,!android,!nox11 freebsd openbsd !windows

package vscode

import (
	"strings"

	"github.com/audrenbdb/goforeground"
)

func switchToWs(name string, pvcs pvcs) bool {

	for _, pvc := range pvcs {
		// https://stackoverflow.com/questions/47189825/golang-how-to-set-window-on-top
		if strings.Contains(pvc.cwd, name) {
			goforeground.Activate(int(pvc.pid))
			return true
		}
	}
	return false
}
