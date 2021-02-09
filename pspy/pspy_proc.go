package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const waitTimeMS = 30

var (
	outputFile *string
	filterUID  *int
	filterCMD  *string
	timeout    *int
)

type process struct {
	cmdline string
	uid     int
	pid     int
}

func (p process) String() string {
	return fmt.Sprintf("PID=%5d UID=%4d CMD=%s", p.pid, p.uid, p.cmdline)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseCommandLineArguments() {
	outputFile = flag.String("w", "-", "output file")
	filterUID = flag.Int("uid", -1, "filter UID")
	filterCMD = flag.String("cmd", "", "filter CMD")
	timeout = flag.Int("timeout", math.MaxInt32, "stop after x seconds")
	flag.Parse()

	if *outputFile != "-" {
		file, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		//FIXME: the log file is never closed. We can't close it here. if we close it in main: does this work if we don't catch the ctrl-c?
		//defer file.Close()
		log.SetOutput(file)
	}
}

// https://stackoverflow.com/questions/30141956/sorting-files-in-numerical-order-in-go/30142239
type byNumericalFilename []os.FileInfo

func (nf byNumericalFilename) Len() int {
	return len(nf)
}

func (nf byNumericalFilename) Swap(i, j int) {
	nf[i], nf[j] = nf[j], nf[i]
}

func (nf byNumericalFilename) Less(i, j int) bool {
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	pathAasInt, err1 := strconv.ParseInt(pathA, 10, 64)
	pathBasInt, err2 := strconv.ParseInt(pathB, 10, 64)

	if err1 != nil || err2 != nil {
		return pathA > pathB
	}
	return pathAasInt > pathBasInt
}

func getProccessesFromProc(pl map[int]process) {

	proc, err := ioutil.ReadDir("/proc")
	check(err)
	sort.Sort(byNumericalFilename(proc))
	for _, entry := range proc {

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			//there are more things in /proc than pids
			continue
		}

		_, pidExists := pl[pid]
		if pidExists {
			// fmt.Println("Process is already known")
			continue
		}

		p := process{cmdline: "???", uid: -1, pid: pid}

		cmdline, err := ioutil.ReadFile("/proc/" + entry.Name() + "/cmdline")
		if len(cmdline) == 0 {
			cmdline, err = ioutil.ReadFile("/proc/" + entry.Name() + "/comm")
		}
		if err == nil {
			cmdlineNice := strings.TrimRight(strings.ReplaceAll(string(cmdline), string(0), " "), "\n")
			if *filterCMD != "" && !strings.Contains(cmdlineNice, *filterCMD) {
				continue
			}
			p.cmdline = cmdlineNice
		}

		status, err := ioutil.ReadFile("/proc/" + entry.Name() + "/status")
		if err == nil {
			regex := regexp.MustCompile(`Uid:\s*(\d+)\s*`)
			uidString := regex.FindStringSubmatch(string(status))[1]
			uid, err := strconv.Atoi(uidString)
			if err != nil {
				log.Print("Error converting uid to int", err)
			}
			if *filterUID != -1 && uid != *filterUID {
				continue
			}
			p.uid = uid
		}

		log.Print(p)
		pl[pid] = p
	}

}

func main() {
	parseCommandLineArguments()
	ticker := time.NewTicker(waitTimeMS * time.Millisecond)
	timeoutTimer := time.NewTicker(time.Duration(*timeout) * time.Second)
	processList := make(map[int]process, 0)
	for {
		select {
		case <-ticker.C:
			getProccessesFromProc(processList)
		case <-timeoutTimer.C:
			log.Print("Time is over")
			os.Exit(0)
		}
	}
}
