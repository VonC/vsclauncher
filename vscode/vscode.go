package vscode

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"vsclauncher/logger"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/hbollon/go-edlib"
	"golang.org/x/sys/windows"
)

type wsfinder struct {
	path        string
	currentPath string
	name        string
	vscode      string
}

type workspace struct {
	fullpath string
	path     string
	name     string
	distance int
}

func newWorkspace(fullpath string, filter string) *workspace {
	path := filepath.Dir(fullpath)
	name := filepath.Base(fullpath)
	name = strings.TrimSuffix(name, ".code-workspace")
	d := edlib.LCSEditDistance(name, filter)
	res := &workspace{
		fullpath: fullpath,
		path:     path,
		name:     name,
		distance: d,
	}
	return res
}

func (w workspace) String() string {
	return w.fullpath
}

type workspaces []*workspace

func (ws workspaces) isUnique() bool {
	return (len(ws) == 1)
}

func (ws workspaces) Len() int           { return len(ws) }
func (ws workspaces) Swap(i, j int)      { ws[i], ws[j] = ws[j], ws[i] }
func (ws workspaces) Less(i, j int) bool { return ws[i].distance < ws[j].distance }

func newWorkspaceFinder(path string, name string) *wsfinder {
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
	wsf := newWorkspaceFinder(path, name)
	ws := wsf.find()
	if ws.isUnique() {
		return ws[0].String()
	}
	wsf.resetCurrentPath()
	ws = wsf.findInParentVScodeFolders()
	if ws.isUnique() {
		return ws[0].String()
	}
	ws = wsf.findInParentGitRoot()
	if ws.isUnique() {
		return ws[0].String()
	}
	wsf.resetCurrentPath()
	sort.Sort(ws)
	/*
		for _, w := range ws {
			fmt.Printf("Workspace (%d) '%s'\n", w.distance, w.fullpath)
		}
	*/
	return getWorkspace(ws).fullpath
}

func (wsf *wsfinder) findInParentVScodeFolders() workspaces {
	if wsf.isGitRoot() {
		return nil
	}
	wsf.currentPath = filepath.Dir(wsf.currentPath)
	logger.Debug("findInParentVScodeFolders '%s'", wsf.currentPath)
	ws := wsf.find()
	if len(ws) > 0 {
		return ws
	}
	return wsf.findInParentVScodeFolders()
}

func (wsf *wsfinder) findInParentGitRoot() workspaces {
	res := make(workspaces, 0)
	wsf.resetCurrentPath()
	p := filepath.Dir(wsf.path)
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		//logger.Debug("file '%s'", name)
		if file.IsDir() {
			wsf.currentPath = filepath.Join(p, name)
			ws := wsf.find()
			res = append(res, ws...)
		}
	}
	return res
}

func (wsf *wsfinder) resetCurrentPath() {
	wsf.currentPath = wsf.path
}

func (wsf *wsfinder) find() workspaces {
	p := wsf.currentPath
	res := make(workspaces, 0)
	// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		//logger.Debug("file '%s'", name)
		if file.IsDir() {
			if name == ".vscode" {
				logger.Debug(".vscode detected")
				wsf.currentPath = filepath.Join(p, name)
				ws := wsf.find()
				wsf.currentPath = p
				res = append(res, ws...)
			} else {
				//logger.Debug("Skip folder '%s'", name)
				continue
			}
		} else if strings.HasSuffix(name, ".code-workspace") {
			w := newWorkspace(filepath.Join(p, name), wsf.name)
			if wsf.hasNoFilter() || w.distance < len(w.name) {
				logger.Debug("Add '%s', distance between '%s' and '%s': '%d'", w, w.name, wsf.name, w.distance)
				res = append(res, w)
			} else {
				logger.Debug("SKIP '%s', distance between '%s' and '%s': '%d'", w, w.name, wsf.name, w.distance)
			}
		} else if name == "_x_" {
			logger.Debug("Skip file '%s'", name)
		}
	}
	return res
}

func (wsf *wsfinder) isGitRoot() bool {
	p := wsf.currentPath
	// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		//logger.Debug("file '%s'", name)
		if file.IsDir() {
			if name == ".git" {
				return true
			}
		}
	}
	return false
}

func (wsf *wsfinder) hasNoFilter() bool {
	return len(wsf.name) == 0
}

func (ws workspaces) names() []string {
	res := make([]string, 0)
	for _, w := range ws {
		res = append(res, w.name)
	}
	return res
}
func (ws workspaces) getWorkspace(name string) *workspace {
	for _, w := range ws {
		if w.name == name {
			return w
		}
	}
	return nil
}

func getWorkspace(ws workspaces) *workspace {

	// https://github.com/charmbracelet/bubbletea/issues/121
	// https://github.com/erikgeiser/coninput/blob/main/example/main.go
	con, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		log.Fatalf("get stdin handle: %s", err)
	}

	var originalConsoleMode uint32

	if runtime.GOOS == "windows" {
		err = windows.GetConsoleMode(con, &originalConsoleMode)
		if err != nil {
			log.Fatalf("get console mode: %s", err)
		}

		defer func() {
			resetErr := windows.SetConsoleMode(con, originalConsoleMode)
			if err == nil && resetErr != nil {
				log.Fatalf("reset console mode: %s", resetErr)
			}
		}()
	}

	sp := selection.New("Chose a VSCode Worskpace to open:",
		selection.Choices(ws.names()))
	sp.PageSize = 3

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	w := ws.getWorkspace(choice.String)
	return w
}
