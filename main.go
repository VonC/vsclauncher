package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"vsclauncher/logger"
	"vsclauncher/vscode"
)

func main() {
	logger.SetLevel(logger.INFO)
	if os.Getenv("VSCDEBUG") != "" {
		logger.SetLevel(logger.DEBUG)
	}
	logger.Debug("vsclauncher (%VSCDEBUG% for more)")
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := ""
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	w := vscode.FindWorkspace(currentWorkingDirectory, name)
	fmt.Printf("Main: workspace found: '%s'\n", w)
	if w == "" {
		m := ""
		if name != "" {
			m = " for name '" + name + "'"
		}
		fmt.Printf("No VSCode workspace found%s\n", m)
	}

	if !vscode.SwitchTo(w) {
		vscode.Launch(w)
		time.Sleep(2 * time.Second)
		vscode.SwitchTo(w)
	}
}
