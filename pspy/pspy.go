package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// "fsnotify"
	"github.com/kmille/fsnotify"
	"golang.org/x/sys/unix"
)

var (
	//command line arguments
	outputFile *string
	filterUID  *int
	filterCMD  *string

	watcher *fsnotify.Watcher
	logFile *os.File
	signals chan os.Signal
)

func parseCommandLineArguments() {
	outputFile = flag.String("w", "-", "output file")
	filterUID = flag.Int("uid", -1, "filter UID")
	filterCMD = flag.String("cmd", "", "filter CMD")
	flag.Parse()

	if *outputFile != "-" {
		logFile, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	}
}

func handleNotifyEvents() {
	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				log.Println("error:", ok)
				return
			}
			// log.Println("New event:", event.Name, event.Op)
			checkForNewProcesses()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error", err)
		case signal := <-signals:
			log.Println("Received signal", signal)
			watcher.Close()
			if logFile != nil {
				logFile.Close()
			}
			os.Exit(0)
		}

	}
}

func main() {
	parseCommandLineArguments()
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// I can't use := here. If so Go would use watcher as a local variable and not as the global variable
	var err error
	watcher, err = fsnotify.NewWatcher(unix.IN_OPEN)
	if err != nil {
		log.Fatal(err)

	}

	go checkForNewProcesses()
	go handleNotifyEvents()

	for _, location := range strings.Split(os.Getenv("PATH"), ":") {
		log.Println("Watching", location)
		err = watcher.Add(location)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Watcher started")
	// wait for ctrl-c
	done := make(chan bool)
	<-done
}
