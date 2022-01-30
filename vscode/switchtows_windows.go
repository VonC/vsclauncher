//go:build windows
// +build windows

package vscode

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

func switchToWs(name string, pvcs pvcs) bool {

	for _, pvc := range pvcs {
		// https://stackoverflow.com/questions/47189825/golang-how-to-set-window-on-top
		if strings.Contains(pvc.cwd, name) {
			pppvc := pvcs.pppid(pvc.pid)
			wh, err := findWindow(int(pppvc.pid), name)
			if err != nil {
				log.Fatalln(err)
			}
			//fmt.Printf("Whnd '%d' for pid '%d', name '%s'\n", wh, pppvc.pid, name)
			wh.ShowWindow(co.SW_SHOWMAXIMIZED)
			setForeground(wh)
			wh.ShowWindow(co.SW_SHOWMAXIMIZED)
			return true
		}
	}

	return false
}

const (
	//Activates and displays the window.
	//If the window is minimized or maximized,
	//the system restores it to its original size and position
	swRestore = 9
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procEnumWindows         = user32.MustFindProc("EnumWindows")
	procShowWindow          = user32.MustFindProc("ShowWindow")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
)

func setForeground(h win.HWND) error {
	procShowWindow.Call(uintptr(h), swRestore)
	procSetForegroundWindow.Call(uintptr(h))
	return nil
}

func findWindow(pid int, name string) (win.HWND, error) {
	//fmt.Printf("findWindow '%d'\n", pid)
	var hwnd win.HWND
	cb := syscall.NewCallback(func(h win.HWND, p uintptr) uintptr {
		t := h.GetWindowText()
		if strings.Contains(t, name+" ") {
			//fmt.Printf("Check Window '%d'", h)
			//fmt.Printf(" title '%s'\n", t)
			hwnd = h
			return 0
		}
		return 1
	})
	enumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("no window with pid %d found", pid)
	}
	return hwnd, nil
}

func enumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
