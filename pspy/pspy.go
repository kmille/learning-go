package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kmille/fsnotify"
	"golang.org/x/sys/unix"
)

var (
	//command line arguments
	outputFile *string
	filterUID  *int
	filterCMD  *string
	debug      *bool

	watcher *fsnotify.Watcher
	logFile *os.File
	signals chan os.Signal
)

func parseCommandLineArguments() {
	outputFile = flag.String("w", "-", "output file")
	filterUID = flag.Int("uid", -1, "filter UID")
	filterCMD = flag.String("cmd", "", "filter CMD")
	debug = flag.Bool("debug", false, "debug print for every event")
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
		case event := <-watcher.Events:
			if *debug {
				log.Println("New event:", event)
			}
			checkForNewProcesses()
		case err := <-watcher.Errors:
			cleanup(fmt.Sprintln("Watcher error", err.Error()), 1)
		case signal := <-signals:
			cleanup(fmt.Sprintln("Received signal", signal), 0)
		}
	}
}

func main() {
	parseCommandLineArguments()
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// I can't use := here. If so Go would use watcher as a local variable and not as the global variable
	var err error
	watcher, err = fsnotify.NewWatcher(unix.IN_OPEN | unix.IN_ACCESS)
	if err != nil {
		log.Fatal(err)
	}

	go checkForNewProcesses()
	go handleNotifyEvents()

	for _, location := range strings.Split(os.Getenv("PATH"), ":") {
		log.Println("Watching", location)
		err = watcher.Add(location)
		if err != nil {
			cleanup(fmt.Sprintln("Error adding a location to the watcher", err), 1)
		}
	}

	log.Println("Watcher started")
	// wait for ctrl-c
	done := make(chan bool)
	<-done
}

func cleanup(message string, returnValue int) {
	log.Println(message)
	watcher.Close()
	if logFile != nil {
		logFile.Close()
	}
	os.Exit(returnValue)
}
