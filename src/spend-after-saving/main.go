/*
Spend After Saving
Function Parameters. +10 XP

"Do not save what is left after spending; instead spend what is left after saving." said Warren Buffett. Inspired by these words Jack decided to save 15% of his monthly salary.
You are given a program that takes salary as input. Complete the function in order to calculate and output the savings.

Sample Input
200

Sample Output
30


@author davorpatech
@since  2021-09-23
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
    salaryAmount, err := AskF64("Input salary amount: ", 200)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }

    fmt.Printf("Your $%v savings are: %v\n", salaryAmount, getSavings(salaryAmount))
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
    if err != nil || v < 0 {
        err = fmt.Errorf("Invalid number: %q", text)
    }
    return
}

func getSavings(salaryAmount float64) float64 {
    return salaryAmount * 15 / 100;
}
