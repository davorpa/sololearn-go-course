/*
Exam results
else if +10 XP

Тhe result of an exam will be determined as follows։
If the score is
88 and above => excellent
40-87   => good
0-39    => fail

You are given a program that takes the score as input.

Task
Complete the code to output the corresponding result (excellent, good, fail) to the console.

Sample Input
78

Sample Output
good

Hint:
Use fmt.Println() to output the result to the console.
Use logical operator && to chain multiple conditions.


@author davorpatech
@since  2021-09-18
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var score float64
    fmt.Print("Enter exam score: ")
    if _, err := fmt.Scanln(&score); err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", "Invalid number.")
        os.Exit(1)
    }
    fmt.Println()

    fmt.Printf("The score mark is: %v\n",
        FormatScore(score))
}

func FormatScore(score float64) (text string) {
    if (score >= 0 && score < 40) {
        text = "fail"
    } else if (score < 88) {
        text = "good"
    } else if (score >= 88) {
        text = "excellent"
    }
    return
}
