/*
Set A World Record
The if else statement +10 XP

The current world record for high jumping is 2.45 meters.
You are given a program that receives as input a number that represents the height of the jump.

Task
Complete the code to:
1. output to the console "new record" if the number is more than 2.45,
2. output to the console "not this time" in other cases.

Sample Input
2.49

Sample Output
new record

Hint:
Note that if the jump height is equal to 2.45 meters, it's not a new world record.


@author davorpatech
@since  2021-09-18
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var height float64
    fmt.Print("Enter jumping height: ")
    if _, err := fmt.Scanln(&height); err != nil || height < 0 {
        fmt.Fprintln(os.Stderr, "ERR!", "Invalid number")
        os.Exit(1)
    }
    fmt.Println()

    fmt.Println(FormatTestNewWorldRecord(
        height))
}

const WR_THRESHOLD float64 = 2.45

func IsNewWorldRecord(height float64) bool {
    return height > WR_THRESHOLD
}

func FormatTestNewWorldRecord(height float64) (text string) {
    if IsNewWorldRecord(height) {
        text = fmt.Sprint("world record")
    } else {
        text = fmt.Sprint("not this time")
    }
    return
}
