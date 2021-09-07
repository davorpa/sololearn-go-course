/**
Office Computers. Math Operators +10 XP

In the office, 2 monitors are connected to each computer.
The first line of the given code takes the number of computers as input.

Task
Complete the code to calculate and output the number of monitors to the console.

Sample Input
10

Sample Output
20

Hint
Since each computer has 2 monitors, simply multiply the count of the computers by 2.

Use the multiplication operator (*).


@author davorpatech
@since  2021-09-07
*/

package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    var text string
    fmt.Print("Computers? ")
    if _, err := fmt.Scanln(&text); err != nil {
        fmt.Fprintln(os.Stderr, "ERR", err)
        os.Exit(1)
    }
    
    computers, err := strconv.Atoi(text)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR", err)
        os.Exit(1)
    }
    
    monitors := applyFactor(computers, MONITORS_PER_COMPUTER)
    fmt.Printf("\nMonitors: %d\n", monitors)
}

const MONITORS_PER_COMPUTER = 2

func applyFactor(n, factor int) int {
    return n * factor
}
