/*
Find the adults. Comparison Math Operators +10 XP

check someone’s age...not just at the bar!

You are given a program that takes the age of the user as input.
Complete the code to check if the user is an adult, and output to the console the corresponding boolean value.

Sample Input
20

Sample Output
true

Hint:
If the user is 18 or older, they’re considered an adult.
fmt.Println(20>18) outputs true.


@author davorpatech
@since  2021-09-12
*/

package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    var text string
    fmt.Print("Input your age: ")
    if _, err := fmt.Scanln(&text); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    
    age, err := Atoui(text)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    
    if isAdult(age) {
        fmt.Println("You are considered an adult")
    } else {
        fmt.Println("Too pretty young to be adult.")
    }
}

func isAdult(age uint) bool {
    return age >= 18
}

func Atoui(s string) (uint, error) {
    ui64, err := strconv.ParseUint(s, 10, 0)
    return uint(ui64), err
}
