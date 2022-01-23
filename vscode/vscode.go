package vscode

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"vsclauncher/logger"
)

type wsfinder struct {
	path        string
	currentPath string
	name        string
	vscode      string
}

type workspace string

func (w workspace) String() string {
	return string(w)
}

type workspaces []workspace

func (ws workspaces) isUnique() bool {
	return (len(ws) == 1)
}

func newWorkspace(path string, name string) *wsfinder {
	wsf := &wsfinder{
		path: filepath.Clean(path),
		name: name,
	}
	wsf.currentPath = path + string(filepath.Separator)
	wsf.vscode = ".vscode" + string(filepath.Separator)
	return wsf
}

func FindWorkspace(path string, name string) string {
	logger.Debug("FindWorkspace in path '%s', name '%s'", path, name)
	wsf := newWorkspace(path, name)
	ws := wsf.find(wsf.currentPath)
	if ws.isUnique() {
		return ws[0].String()
	}
	return ""
}

func (wsf *wsfinder) find(p string) workspaces {
	res := make(workspaces, 0)
	// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		logger.Debug("file '%s'", name)
		if file.IsDir() {
			if name == ".vscode" {
				logger.Debug(".vscode detected")
				ws := wsf.find(filepath.Join(p, name))
				res = append(res, ws...)
			} else {
				logger.Debug("Skip folder '%s'", name)
				continue
			}
		} else if strings.HasSuffix(name, ".code-workspace") {
			res = append(res, workspace(strings.TrimSuffix(name, ".code-workspace")))
		} else {
			logger.Debug("Skip file '%s'", name)
		}
	}
	return res
}
