/*
Number of Ones +10 XP

Task:
Count the number of ones in the binary representation of a given integer.

Input Format:
An integer.

Output Format: 
An integer representing the count of ones in the binary representation of the input.

Sample Input:
9

Sample Output:
2

Explanation: 
The binary representation of 9 is 1001, which includes 2 ones.


@author davorpatech
@since 2021-08-16
*/

package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func init() {
    log.SetFlags(log.Ltime|log.Lmicroseconds|log.LUTC)
    log.SetOutput(os.Stdout)
}

func main() {
    var value string
    fmt.Print("Input a number: ")
    if _,err := fmt.Scan(&value); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    
    // as unsigned integer
    i64, err := strconv.ParseUint(value, 10, 64)
    if err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    
    // to base-2
    bin := strconv.FormatUint(i64, 2)
    
    log.Print(i64)
    log.Print(bin)
    
    // count ones
    ones := strings.Count(bin, "1")
    fmt.Println("Number of ones:", ones)
}
