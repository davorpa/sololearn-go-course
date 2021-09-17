/*
Vacation Month
The if Statement +10 XP

You are planning a vacation in August.
You are given a program that takes the month as input.

Task
Complete the code to output "vacation", if the given month is August. Don't output anything otherwise.

Sample Input
August

Sample Output
vacation

Hint:
Handle the required condition using an if statement.
Use fmt.Println() for the output.


@author davorpatech
@since  2021-09-16
*/

package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    var month string
    fmt.Print("Input month: ")
    fmt.Scanln(&month)
    fmt.Println()

    if matches, err := IsVacationalMonth(month); err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    } else if matches {
        fmt.Println("vacational")
    } else {
        fmt.Println("no vacational")
    }
}

var months = map[string]byte{
    "January": 0,
    "February": 0,
    "March": 0,
    "April": 0,
    "May": 0,
    "June": 0,
    "July": 0,
    "August": 1,
    "September": 0,
    "October": 0,
    "November": 0,
    "December": 0,
    }

func IsVacationalMonth(month string) (matches bool, err error) {
    found := false
    for k, v := range months {
        if equalsCaseInsensitive(k, month) {
            matches = v == 1
            found = true
        }
    }
    if !found {
        err = fmt.Errorf("Invalid month: %s", month)
    }
    return
}

func equalsCaseInsensitive(a, b string) bool {
    return strings.ToLower(a) == strings.ToLower(b)
}
