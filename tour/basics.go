package main

import (
    "fmt"
    "math/rand"
    "math"
)


var c, python, java bool
var i, j int = 1,2
const Pi = 3.14

func add(x int, y int) int {
    // can be written as func add(x, y int) int
    return x + y
}

func swap(a, b string) (string, string) {
    return b, a
}

func split(sum int) (x, y int) {
    y = sum / 10
    x = sum % 10
    // this will return x,y or 0 if x or y is noet defined
    return
}

func data_types() {
    var MaxInt uint64 = 1<<64 -1
    var c byte = 14
    fmt.Println("This is max int for 64 bit:", MaxInt)
    fmt.Printf("This is one byte %d bin=%b type=%T\n", c, c, c)

 }

func main() {
    fmt.Println("My first random number is", rand.Intn(10))
    fmt.Println("My second random number is", rand.Intn(10))
    s := math.Sqrt(7)
    fmt.Printf("%T %v %g\n", s, s, s)
    fmt.Println(add(1, 4))
    a, b := swap("world", "hello,")
    fmt.Println(a, b)
    fmt.Println(split(19))
    c = true
    fmt.Println("global booleans", c, python, java)
    fmt.Println("global initialized variables", i, j)
    k := 3 // instead of var k = 3
    fmt.Println("Just the variable k", k)
    data_types()
    fmt.Println("This is pi:", Pi)
}
