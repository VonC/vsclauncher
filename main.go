package main

import (
	"log"
	"os"
	"vsclauncher/logger"
	"vsclauncher/vscode"
)

func main() {
	logger.SetLevel("debug")
	logger.Debug("vsclauncher")
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	vscode.FindWorkspace(currentWorkingDirectory)
}
