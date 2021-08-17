/*
Multiples +10 XP


You need to calculate the sum of all the multiples of 3 or 5 below a given number.

Task: 
Given an integer number, output the sum of all the multiples of 3 and 5 below that number. 
If a number is a multiple of both, 3 and 5, it should appear in the sum only once.

Input Format: 
An integer.

Output Format: 
An integer, representing the sum of all the multiples of 3 and 5 below the given input.

Sample Input: 
10

Sample Output:
23

ℹ Explanation: 
The numbers below 10 that are multiples of 3 or 5 are 3, 5, 6 and 9.
The sum of these numbers is 3+5+6+9=23


@author davorpatech
@since  2021-08-17
*/
package main

import (
    "fmt"
    "os"
)

func main() {
    var n uint
    fmt.Print("Input a number: ")
    if _, err := fmt.Scan(&n); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    
    var sum uint
    for i := uint(1); i < n; i++ {
        if i % 3 == 0 || i % 5 == 0 {
            sum += i
        }
    }
    fmt.Println("Sum of 3|5 multiples below:", sum)
}
