package main

import (
    "fmt"
    "runtime"
)


func shutdown(c int) {
    fmt.Println("I'm run at the end as number: ", c)
}

func main() {

    // implemented as a stack
    defer shutdown(1)
    defer shutdown(2)

    for i :=0; i< 10; i++ {
        fmt.Printf("%d ", i)
    }
    fmt.Printf("\n")

    var sum = 1
    for sum < 10 {
        sum += sum
    }
    fmt.Println(sum)

    if 100 > 4 {
        fmt.Println("Yes!")
    } else {
        fmt.Println("Nooooooo!")
    }

    var os = runtime.GOOS
    switch os {
    case "linux":
        fmt.Println("The best")
        // here is an explicit break
    default:
        fmt.Println("default")

    }

}

