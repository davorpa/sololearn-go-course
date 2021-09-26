/*
Average Of Three
The return Statement. +10 XP

You are given a program that takes 3 numbers as input.
Complete the given function to calculate the average of those 3 numbers, assign it to the given variable, and output it.

Sample Input
3
6
9

Sample Output
6

Hint:
Use return statement to return the calculated value and use it through the program.

@author davorpatech
@since  2021-09-24
*/

package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    nums, err := ScanNums(3, 0)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }
    
    fmt.Printf("%d-Average is: %f",
        len(nums), Avg(nums...))
}

func ScanNums(min, max uint) (nums []float64, err error) {
    if max > 0 {
        nums = make([]float64, 0, max)
    }
    for i := uint(0); max == 0 || i < max; {
        var (
            text string
            n float64
        )
        i++
        fmt.Printf("Number %d (blank=end): ", i)
        fmt.Scanln(&text)
        fmt.Println()
        
        text = strings.TrimSpace(text)
        if len(text) > 0 {
            n, err = strconv.ParseFloat(text, 64)
            if err != nil {
                if _, ok := err.(*strconv.NumError); ok {
                    err = fmt.Errorf("Invalid number %d: %q", i, text)
                }
                return
            }
            nums = append(nums, n)
            continue
        }
        if min >= i {
            err = fmt.Errorf("Input at least %d numbers.", min)
        }
        break
    }
    return
}

func Avg(v ...float64) (m float64) {
    l := len(v)
    if l == 0 {
        return
    }
    for _, n := range v {
        m += n
    }
    m = m / float64(l)
    return
}
