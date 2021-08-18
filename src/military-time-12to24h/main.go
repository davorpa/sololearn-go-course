/*
Military Time +10 XP

You want to convert the time from a 12 hour clock to a 24 hour clock. If you are given the time on a 12 hour clock, you should output the time as it would appear on a 24 hour clock.  

Task:  
Determine if the time you are given is AM or PM, then convert that value to the way that it would appear on a 24 hour clock.

Input Format: 
A string that includes the time, then a space and the indicator for AM or PM.

Output Format: 
A string that includes the time in a 24 hour format (XX:XX)

Sample Input: 
1:15 PM

Sample Output: 
13:15

Explanation:
1:00 PM on a 12 hour clock is equivalent to 13:00 on a 24 hour clock.


@author davorpatech
@since  2021-08-18
*/

package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
    )

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)
    
    fmt.Print("Input time (h:mm AM|PM): ")
    if scanner.Scan() {
        input := scanner.Text()
        
        mt, err := MilitaryTime(input)
        if err != nil {
            fmt.Fprint(os.Stderr, err)
            os.Exit(1)
        }
        fmt.Printf("24h time: %s", mt)
        return
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
}


func MilitaryTime(st string) (mt string, err error) {
    var dt time.Time
    // case insensitive support
    st = strings.ToUpper(st)
    // from 12h format
    dt, err = time.Parse("3:04 PM", st)
    if err == nil {
        // to 24h format
        mt = dt.Format("15:04")
    }
    return
}
