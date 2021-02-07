package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func channelExample() {
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c /// reading from the channel, block
	fmt.Println(x, y, x+y)
}

func firstConcurrencyExample() {
	go say("world")
	say("hello")
}

func channelExampleWithBufferLength() {
	// channel with a buffer length. Writes block if you write more then length
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func fibonacci(c chan int) {
	x, y := 0, 1
	for i := 0; i < cap(c); i++ {
		c <- x
		x, y = y, x+y
	}
	//this will end the loop
	// close is only needed if the other party needs to know if the channel is closed
	close(c)
}

func fibonacciExample() {
	c := make(chan int, 10)
	go fibonacci(c)
	//this will block until we have data
	for i := range c {
		fmt.Println(i)
	}

}

func fibonacci2(c chan int, quit chan int) {
	x, y := 0, 1
	for {
		select {
		//write int x into c
		case c <- x:
			x, y = y, x+y
		// read from quit
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func channelExampleSwitch() {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			//read int from c
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci2(c, quit)
}

func channelExampleSwitchDefault() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func main() {
	//firstConcurrencyExample()
	//channelExample()
	//channelExampleWithBufferLength()
	//fibonacciExample()
	//channelExampleSwitch()
	channelExampleSwitchDefault()
	//TODO : sync.Mutex

}
