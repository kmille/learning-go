package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestGetPIDs(t *testing.T) {
	sleepProcess, err := runNewProcess(t)
	if err != nil {
		t.Log("Could not start sleep", err)
	}
	pids, err := getPIDs()
	if err != nil {
		cleanupTest(t, sleepProcess)
		t.Error("Could not get the PIDs", err)
	}
	if !contains(t, pids, sleepProcess.Pid) {
		cleanupTest(t, sleepProcess)
		t.Errorf("Could not find pid %d in all pids (%d)", sleepProcess.Pid, pids)
	}
	cleanupTest(t, sleepProcess)
	t.Log("Done testing getPIDs")
}

func TestGetProcessCmdline(t *testing.T) {
	// this is shitty. If not called then *filterCMD points to garbage
	// we also can't call parseCommandLineArguments twice
	parseCommandLineArguments()
	sleepProcess, err := runNewProcess(t)
	if err != nil {
		t.Log("Could not start sleep", err)
	}
	cmdline, err := getProcessCmdline(sleepProcess.Pid)
	if err != nil {
		cleanupTest(t, sleepProcess)
		t.Error("Could not get the PIDs", err)
	}
	if cmdline != "/bin/sleep 1" {
		cleanupTest(t, sleepProcess)
		t.Errorf("cmdline is not what we expected: %q", cmdline)
	}
	cleanupTest(t, sleepProcess)
	t.Log("Done testing TestGetProcessCmdline")
}

func TestGetProcessUID(t *testing.T) {
	sleepProcess, err := runNewProcess(t)
	if err != nil {
		t.Log("Could not start sleep", err)
	}
	uid, err := getProcessUID(sleepProcess.Pid)
	if err != nil {
		cleanupTest(t, sleepProcess)
		t.Error("Could not get the uid", err)
	}
	if uid != os.Getuid() {
		cleanupTest(t, sleepProcess)
		t.Errorf("uid is not what we expected: %d", uid)
	}
	cleanupTest(t, sleepProcess)
	t.Log("Done testing getProcessUID")
}

func runNewProcess(t *testing.T) (*os.Process, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("/bin/sleep", "1")
	cmd.Dir = cwd
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	t.Log("running /bin/sleep in the background with pid", cmd.Process.Pid)
	return cmd.Process, nil
}

func contains(t *testing.T, ints []int, element int) bool {
	for _, i := range ints {
		if element == i {
			return true
		}
	}
	return false
}

func cleanupTest(t *testing.T, sleepProcess *os.Process) {
	sleepProcess.Wait()
	sleepProcess.Release()
}
