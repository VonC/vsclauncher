package vscode

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"vsclauncher/internal/syscall"
	"vsclauncher/logger"
)

func Launch(w string) {
	ws := newWorkspace(w, "")
	logger.Debug("Launch '%s'\n", w)
	d := os.Getenv("vscodei")
	if d == "" {
		d = filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs\\Microsoft VS Code\\bin\\code.cmd")
		if _, err := os.Stat(d); os.IsNotExist(err) {
			d = filepath.Join(os.Getenv("ProgramFiles"), "Microsoft VS Code\\bin\\code.cmd")
		}
	} else {
		d = filepath.Join(d, "code.exe")
	}
	// https://stackoverflow.com/questions/6376113/how-do-i-use-spaces-in-the-command-prompt
	c := fmt.Sprintf(`"cd "%s" && "%s" "%s""`, ws.path, d, w)
	//fmt.Println(c)

	giterr, _, err := syscall.ExecCmd(c)
	if err != nil || giterr.String() != "" {
		log.Fatalf("Unable to launch VSCode workspace '%s': '%s' (%+v)", w, giterr.String(), err)
	}
}
