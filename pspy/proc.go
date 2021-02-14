package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var processList = make(map[int]process)

type process struct {
	cmdline string
	uid     int
	pid     int
}

func (p process) String() string {
	return fmt.Sprintf("PID=%-5d UID=%-5d CMD=%s", p.pid, p.uid, p.cmdline)
}

func getProcessCmdline(pid int) (string, error) {
	// TODO: error handling
	cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil || (len(cmdline) == 0) {
		cmdline, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	}
	// log.Println(err)
	if err == nil {
		cmdlineNice := strings.TrimRight(strings.ReplaceAll(string(cmdline), string(0), " "), "\n")
		if *filterCMD != "" && !strings.Contains(cmdlineNice, *filterCMD) {
			return "", errors.New(fmt.Sprintf("cmdline filter: '%s'", cmdlineNice))
		}
		return cmdlineNice, nil
	}

	return "", errors.New("Could not extract cmdline")
}

func printProcessInfos(pid int) {
	p := process{cmdline: "???", uid: -1, pid: pid}

	cmdline, err := getProcessCmdline(pid)
	if err != nil {
		log.Printf("Could not get cmdline: %s\n", err)
	}
	p.cmdline = cmdline

	statInfo := syscall.Stat_t{}
	err = syscall.Lstat(fmt.Sprintf("/proc/%d/", pid), &statInfo)
	if err == nil {
		p.uid = int(statInfo.Uid)
	}
	if *filterUID != -1 && p.uid != *filterUID {
		return
	}
	log.Print(p)
	processList[pid] = p
}

func checkForNewProcesses() {

	pids, err := getPIDs()
	if err != nil {
		log.Println("Error reading from /proc")
	}
	// for i := len(pids) - 1; i >= 0; i-- {
	for i := 0; i < len(pids); i++ {
		_, exist := processList[pids[i]]
		if !exist {
			printProcessInfos(pids[i])

		}
	}
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
