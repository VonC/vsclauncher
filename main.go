package main

import (
	"fmt"
	"log"
	"os"
	"vsclauncher/logger"
	"vsclauncher/vscode"
)

func main() {
	logger.SetLevel(logger.DEBUG)
	logger.Debug("vsclauncher")
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := ""
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	w := vscode.FindWorkspace(currentWorkingDirectory, name)
	logger.Debug("Main: workspace found: '%s'", w)
	if w == "" {
		m := ""
		if name != "" {
			m = " for name '" + name + "'"
		}
		fmt.Printf("No VSCode workspace found%s\n", m)
	}
}
