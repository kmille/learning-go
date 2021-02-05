package main

import (
    "fmt"
    "strings"
    "math"
)


type ExampleStruct struct {
    X int
    Y int
}

func tic_tac_toe() {
    board := [][]string{
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
        []string{"_", "_", "_"},
    }


    board[0][0] = "X"
    board[2][2] = "O"
    board[1][2] = "X"
    board[1][0] = "O"
    board[0][2] = "X"


    for i := 0; i < len(board); i++ {
        fmt.Printf("%s\n", strings.Join(board[i], " "))
    }
}


func compute(fn func(float64, float64) float64) float64 {
    return fn(10, 2)
}

func adder() func(int) int {
    //this is called closure.?
     sum := 0
     return func(x int) int {
         sum += x
         return sum
     }
}

func create_a_map() {
    type Vertex struct {
        Lat, Long float64
    }
    //var m map[string]Vertex
    // dictionary: key=string value=Vertex
    m := make(map[string]Vertex)
    m["Bell Labs"] = Vertex{ 123.123, 345.345 }
    fmt.Println(m)
    fmt.Println(m["Bell Labs"])


    var m2 = map[string]Vertex{
        "Bell Labs": Vertex{1, 2},
        "Google": {2, 3},
    }
    fmt.Println(m2)

    m2["robo6"] = Vertex{ 45435, 123123}
    fmt.Println(m2)

    delete(m2, "robo6")
    fmt.Println(m2)

    v, ok := m2["robo6"]
    fmt.Println("The value:", v, "Present?", ok)
}

func append_to_a_slice() {
    var s []int
    printSlice(s)

    s = append(s, 0, 1, 2)
    printSlice(s)

}

func printSlice(s []int) {
    fmt.Printf("len=%d, cap=%d %v\n", len(s), cap(s), s)
}


func ranger() {
    r := []int{10, 22, 43, 64, 51, 62}
    for i, v := range r {
        fmt.Printf("index=%d, value=%d\n", i, v)
    }

    for _, v := range r {
        fmt.Printf("value=%d\n", v)
    }
    for i:= range r {
        fmt.Printf("index=%d\n", i)
    }
}

func main() {
    var i = 43
    var p = &i
    fmt.Println("Adresses:", &i, p)
    *p = 200
    fmt.Println("Values:", i, *p)

    fmt.Println(ExampleStruct{1, 2})
    fmt.Println("empty struct", ExampleStruct{})
    fmt.Println(ExampleStruct{Y: 15})

    var es = ExampleStruct{1, 2}
    es.X = 100
    fmt.Println(es)


    var a[2]string
    a[0] = "hello"
    a[1] = "world"
    fmt.Println(a[0], a[1])
    fmt.Println(a)

    primes := [6]int{2, 3, 5, 7, 11}
    fmt.Println(primes)

    var part = primes[1:4]
    fmt.Println(part)

    // part still points to primes!
    part[1] = 99
    fmt.Println(primes)


    s := []int{2, 3, 5, 7, 11, 13}
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
    s = s[2:]
    fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)

    var s2 []int
    if s2 == nil {
        fmt.Printf("len=%d cap=%d %v\n", len(s2), cap(s2), s2)
    }

    a2 := make([]int, 5)
    fmt.Println(a2)

    tic_tac_toe()
    append_to_a_slice()
    ranger()
    create_a_map()

    inline_function := func(x, y float64) float64 {
        return math.Sqrt(x*x + y*y)
    }
    fmt.Println(compute(inline_function))


    pos, neg := adder(), adder()
    for i := 0; i < 10; i++ {
        fmt.Println(
            pos(i),
            neg(-2*i),
        )
    }

    fmt.Println(fibonacci(10))
    f := fibonacci2()
	for i := 0; i < 3; i++ {
		fmt.Println(f())
	}

}

func fibonacci(length int) []int {
    fib := []int{}
    for i := 0; i < length; i++ {
        if i == 0 {
            fib = append(fib, 0)
        } else if i == 1 {
            fib = append(fib, 1)
        } else {
            fib = append(fib, fib[len(fib) -2] + fib[len(fib) - 1])
        }
    }
    return fib
}

func fibonacci2() func() int {
    fib := 0
    return func() int {
        fmt.Println("Calculating the", fib, "st fib number")
        first, second := 0, 1
        for i := 0; i < fib; i++ {
            first, second = second, first + second
        }
        fib++
        return first
    }
}
