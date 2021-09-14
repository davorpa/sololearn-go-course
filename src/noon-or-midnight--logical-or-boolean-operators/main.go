/*
Noon or Midnight +10 XP
(Logical or boolean operators 2)

Time flies when you’re having fun.
Given a clock that measures 24 hours in a day, write a program that takes the hour as input. If the hour is in the range of 0 to 12, output am to the console, and output pm if it's not.

Sample Input
13

Sample Output
pm

Hint:
Assume the input number is positive and less than or equal to 24.


@author davorpatech
@since  2021-09-14
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var hour uint8
    fmt.Print("Input an hour (0-24): ")
    if _, err := fmt.Scanln(&hour); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    fmt.Println(ToMeridiemFormat(hour))
}

func ToMeridiemFormat(hour uint8) (s string) {
    h, t, zs := ToMeridiem(hour)
    var tl string
    if zs {
        tl = "post meridiem"
    } else {
        tl = "ante meridiem"
    }
    s = fmt.Sprintf("%2d%s is %s", h, t, tl)
    return
}

func ToMeridiem(hour uint8) (
        uint8, string, bool) {
    if hour > 24 {
        panic(fmt.Sprintf("Hour must be positive and less or equal 24. Value: %v",
            hour))
    }
    status := hour > 12
    var text string
    if status {
        text = "pm"
        hour -= 12
    } else {
        text = "am"
    }
    
    return hour, text, status
}
