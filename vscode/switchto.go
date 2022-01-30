package vscode

import (
	"log"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type pvc struct {
	pid  int32
	ppid int32
	cwd  string
}

type pvcs map[int32]*pvc

func (pvcs pvcs) add(pid, ppid int32, cwd string) {
	pvc := &pvc{
		pid:  pid,
		ppid: ppid,
		cwd:  cwd,
	}
	pvcs[pvc.pid] = pvc
}

func (pvcs pvcs) pppid(pid int32) *pvc {
	pvc := pvcs[pid]
	//fmt.Println(valast.String(pvc))
	ppvc := pvcs[pvc.ppid]
	//fmt.Println(valast.String(ppvc))
	pppvc := pvcs[ppvc.ppid]
	//fmt.Println(valast.String(pppvc))
	return pppvc
}

func SwitchTo(w string) bool {

	pvcs := make(pvcs)
	ps, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range ps {
		exeName, err := p.Exe()
		if err != nil {
			continue
		}
		cwd, err := p.Cwd()
		if err != nil {
			continue
		}
		ppid, err := p.Ppid()
		if err != nil {
			continue
		}
		if strings.HasSuffix(exeName, "Code.exe") {
			pvcs.add(p.Pid, ppid, cwd)
			//fmt.Printf("Pid %d => '%s'\n", p.Pid, cwd)
		}
	}
	ws := newWorkspace(w, "")
	return switchToWs(ws.name, pvcs)
}
