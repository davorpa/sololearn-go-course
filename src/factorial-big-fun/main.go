/*
Factorial Fun
The for loop 2. +10 XP

A number's factorial is the product of all positive integers less than or equal to the number.
Write a program that takes a number as input and outputs its factorial to the console.

Sample Input
5

Sample Output
120

Explanation
5*4*3*2*1 = 120


@author davorpatech
@since  2021-09-20
*/

package main

import (
    "math/big"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

func main() {
    var text string
    fmt.Print("Enter number: ")
    _, err := fmt.Scanln(&text)
    if errors.Is(err, io.EOF) {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }
    text = DefaultIfBlank(text, "5")
    
    number, err := strconv.ParseUint(text, 10, 0)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", "Invalid number:", text)
        os.Exit(1)
    }
    fmt.Println()

    fmt.Printf("%v! = %v\n",
        number,
        Factorial(number))
}

func DefaultIfBlank(s, fallback string) string {
    if IsBlank(s) {
        s = fallback
    }
    return s
}

func IsBlank(s string) bool {
    s = strings.TrimSpace(s)
    return len(s) == 0
}

func Factorial(n uint64) *big.Int {
    f := big.NewInt(1)
    for x := int64(n); x > 1; x-- {
        f.Mul(f, big.NewInt(x))
    }
    return f
}
