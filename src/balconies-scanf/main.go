/*
Balconies +10 XP

You are trying to determine which of two apartments has a larger balcony. Both balconies are rectangles, and you have the length and width, but you need the area.

Task 
Evaluate the area of two different balconies and determine which one is bigger.

Input Format 
Your inputs are two strings where the measurements for height and width are separated by a comma. The first one represents apartment A, the second represents apartment B.

Output Format: 
A string that says whether apartment A or apartment B has a larger balcony.

Sample Input 
'5,5'
'2,10'

Sample Output 
Apartment A

Explanation 
Since the area of apartment A's balcony is 25 and the area of apartment B's balcony is 20, Apartment A is the correct answer.


@author davorpatech
@since 2021-08-16
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var areaA, areaB uint
    
    if length, width, err := input("A"); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    } else {
        areaA = length * width
    }
    
    if length, width, err := input("B"); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    } else {
        areaB = length * width
    }
    
    if areaA > areaB {
        fmt.Print("Apartment A")
    } else {
        fmt.Print("Apartment B")
    }
}

func input(label string) (length, width uint, err error) {
    fmt.Printf("Enter (length,width) of Balcony %s: ", label)
    _, err = fmt.Scanf("%d,%d", &length, &width)
    return
}
