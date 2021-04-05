package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func speak(text string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("espeak \"%s\" --stdout | paplay --volume 65000", text))
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}
func main() {
	var duration, counter int

	flag.IntVar(&duration, "d", 10, "sleep time")
	flag.Parse()

	fmt.Println("Lauch for life!")
	speak("Let's go du Lauch!")

	for {
		c := fmt.Sprintf("%d", 10*(counter%6))
		go speak(c)
		counter++
		time.Sleep(time.Duration(duration) * time.Second)
	}
}
