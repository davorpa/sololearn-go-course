/*
Repeat and Learn Code! Code! Code!
The for loop 2. +10 XP

"Repetition is the mother of learning, the father of action, which makes it the architect of accomplishment." - Zig Ziglar.
Inspired by these words, let's create a little program that will output an expression which is given as input, 3 times.

Task
Complete the code to output the given expression 3 times.

Sample Input
Learning is fun!

Sample Output
Learning is fun!
Learning is fun!
Learning is fun!

Hint:
Use the for loop to run the required part of the code as many times as needed.


@author davorpatech
@since  2021-09-19
*/

package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)
    
    fmt.Print("Enter text to repeat: ")
    var expression string
    if scanner.Scan() {
        expression = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    expression = DefaultIfBlank(expression,
        "Learning is fun!")

    Rprintln(3, expression)
}

func DefaultIfBlank(s string, fallback string) string {
    if IsBlank(s) {
        s = fallback
    }
    return s
}

func IsBlank(s string) bool {
    s = strings.TrimSpace(s)
    return len(s) == 0
}

func Rprintln(times uint, text string) {
    RFprintln(times, os.Stdout, text)
}

func RFprintln(times uint, writer io.Writer, text string) {
    for i := uint(0); i < times; i++ {
        fmt.Fprintln(writer, text)
    }
}
