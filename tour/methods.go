package main

import (
    "fmt"
    "math"
)


type Vertex struct {
    X, Y float64
}

func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
    // value receiver
    // return a new float
    if f < 0 {
        return float64(-f)
    }
    return float64(f)
}

func (v *Vertex) Scale(f float64) {
    // pointer receiver
    // change values of the v object
    v.X = v.X * f
    v.Y = v.Y * f
}

// BEGIN: Abs ans Scale implemented as classic functions 
func Scale(v *Vertex, f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func Abs(v Vertex) float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
// END: Abs ans Scale implemented as classic functions 


type MyInterface interface {
    My()
}

type Tiae struct {
    Susi string
}


func (t Tiae) My() {
    fmt.Println(t.Susi)
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())

    f := MyFloat(-23)
    fmt.Println(f.Abs())

    v.Scale(10)
    fmt.Println(v.Abs())

    v = Vertex{3, 4}
    Scale(&v, 10)
    fmt.Println(v.Abs())

    var blubb MyInterface = Tiae{"hello"}
    blubb.My()

}
