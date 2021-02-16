package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type process struct {
	cmdline string
	uid     int
	pid     int
}

var processList = make(map[int]process)

func (p process) String() string {
	return fmt.Sprintf("PID=%-5d UID=%-5d CMD=%s", p.pid, p.uid, p.cmdline)
}

func checkForNewProcesses() {

	pids, err := getPIDs()
	if err != nil {
		log.Println("Error reading from /proc", err)
		return
	}
	// iterating backwards is faster but the first print of the processes currently running would be in the wrong order
	// for i := len(pids) - 1; i >= 0; i-- {
	for i := 0; i < len(pids); i++ {
		_, exist := processList[pids[i]]
		if !exist {
			printProcessInfos(pids[i])
		}
	}
}

func printProcessInfos(pid int) {
	p := process{cmdline: "???", uid: -1, pid: pid}

	cmdline, err := getProcessCmdline(pid)
	if err != nil {
		// error or cmdline filtered
		return
	}
	p.cmdline = cmdline

	uid, err := getProcessUID(pid)
	if err != nil {
		// error or uid filtered
		return
	}
	p.uid = uid

	log.Print(p)
	processList[pid] = p
}

func getProcessUID(pid int) (int, error) {
	// return only an error if we filter the UID
	statInfo := syscall.Stat_t{}
	err := syscall.Lstat(fmt.Sprintf("/proc/%d/", pid), &statInfo)
	if err != nil {
		return -1, nil
	}
	uid := int(statInfo.Uid)
	if *filterUID != -1 && uid != *filterUID {
		return -1, fmt.Errorf("Filtered out UID %d", uid)
	}
	return uid, nil
}

func getProcessCmdline(pid int) (string, error) {
	// return only an error if we filter the process
	cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil || (len(cmdline) == 0) {
		cmdline, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	}
	if err != nil {
		return fmt.Sprintf("??? (error: %q)", err), nil
	}
	cmdlineNice := strings.ReplaceAll(string(cmdline), "\000", " ")
	cmdlineNice = strings.TrimRight(cmdlineNice, "\n")
	cmdlineNice = strings.TrimRight(cmdlineNice, " ")
	if *filterCMD != "" {
		if !strings.Contains(cmdlineNice, *filterCMD) {
			return "", fmt.Errorf("cmdline filter: '%s'", cmdlineNice)
		}
	}
	return cmdlineNice, nil

}

func getPIDs() ([]int, error) {
	fd, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	names, err := fd.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var pids = []int{}
	for _, name := range names {
		pid, err := strconv.Atoi(name)
		if err == nil {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}
