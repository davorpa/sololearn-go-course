/*
Cuboid Volume. +10 XP
Introducing Objects

The given structure represents a cuboid (e.g. a rectangular prism) that contains the properties of length, width, and height.
Complete the program to calculate and output the volume of given cuboid to the console.

Hint
To calculate the volume of cuboid use length*width*height formula.


@author davorpatech
@since  2021-09-25
*/

package main

import (
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

func main() {
    fields := [...]string{"length", "width", "height"}
    measures := [...]float64{25, 50, 200}
    for i, f := range fields {
        v, err := AskF64("Input " + f + ": ", measures[i])
        if err != nil {
            fmt.Fprintln(os.Stderr, "ERR!", err)
            os.Exit(1)
        }
        measures[i] = v
    }

    var prism Prism
    var err error
    prism, err = Cuboid(
        measures[0], measures[1], measures[2])
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }
    fmt.Println()
    
    fmt.Printf("%v volume: %f", prism, prism.GetVolume())
}

func AskF64(prompt string, fallback float64) (v float64, err error) {
    if len(prompt) > 0 {
        fmt.Print(prompt)
    }
    var text string
    _, err = fmt.Scan(&text)
    fmt.Println()
    if err != nil && !errors.Is(err, io.EOF) {
        return
    }

    if len(strings.TrimSpace(text)) == 0 {
        return fallback, nil
    }

    v, err = strconv.ParseFloat(text, 64)
    if err != nil {
        err = fmt.Errorf("Invalid number: %q", text)
    }
    return
}

type Prism interface {
    GetVolume() float64
}

type cuboid struct {
    length, width, height float64
}

func Cuboid(length, width, height float64) (c cuboid, err error) {
    defer panicCatcher(&err)
    c = cuboid{}
    c.length = ensureGt0(length, "Cuboid length")
    c.width = ensureGt0(width, "Cuboid width")
    c.height = ensureGt0(height, "Cuboid height")
    return
}

func (c cuboid) String() string {
    return fmt.Sprintf("Cuboid(%f × %f × %f)", c.length, c.width, c.height)
}

func (c cuboid) GetVolume() float64 {
    return c.length * c.width * c.height
}

func ensureGt0(v float64, lbl string) float64 {
    if v < 0 {
        panic(fmt.Sprintf("%s must be a positive number: %f", lbl, v))
    }
    return v
}

// transform panics to errors
func panicCatcher(err *error) {
    if r := recover(); r != nil {
        e, ok := r.(error)
        if !ok {
            e = errors.New(fmt.Sprint(r))
        }
        *err = e
    }
}
