package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const wait_time_ms = 50
type process struct {
	cmdline string
	uid     int
	pid     int
}

func (p process) String() string {
	return fmt.Sprintf("PID=%d UID=%d CMD=%s", p.pid, p.uid, p.cmdline)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getProccessesFromProc(pl map[int]process) {

	proc, err := ioutil.ReadDir("/proc")
	check(err)
	for _, entry := range proc {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		_, pidExists := pl[pid]
		if pidExists {
			// fmt.Println("Process is already known")
			continue
		}

		cmdline, err := ioutil.ReadFile("/proc/" + entry.Name() + "/cmdline")
		if len(cmdline) == 0 {
			cmdline, err = ioutil.ReadFile("/proc/" + entry.Name() + "/comm")
			// fmt.Println("This comes from /comm")
		}

		status, err := ioutil.ReadFile("/proc/" + entry.Name() + "/status")
		if err != nil {
			fmt.Println("Too late for this one:", pid)
			continue
		}
		regex := regexp.MustCompile(`Uid:\s*(\d+)\s*`)
		uidString := regex.FindStringSubmatch(string(status))[1]
		uid, err := strconv.Atoi(uidString)
		if err != nil {
			fmt.Println("Error converting uid to int", err)
		}
		cmdlineNice := strings.TrimRight(strings.ReplaceAll(string(cmdline), string(0), " "), "\n")
		p := process{cmdline: cmdlineNice, uid: uid, pid: pid}
		fmt.Println("New process:", p)
		pl[pid] = p
	}

}

func main() {
	ticker := time.NewTicker(wait_time_ms * time.Millisecond)
	processList := make(map[int]process, 0)
	for {
		select {
		case <-ticker.C:
			getProccessesFromProc(processList)
		}
	}
}
