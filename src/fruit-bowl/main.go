/*
Fruit Bowl +10 XP

You have a bowl on your counter with an even number of pieces of fruit in it. Half of them are bananas, and the other half are apples. You need 3 apples to make a pie. 

Task 
Your task is to evaluate the total number of pies that you can make with the apples that are in your bowl given to total amount of fruit in the bowl.

Input Format
An integer that represents the total amount of fruit in the bowl.

Output Format
An integer representing the total number of whole apple pies that you can make.

Sample Input
26 

Sample Output 
4


@author davorpatech
@since  2021-08-19
*/
package main

import (
    "fmt"
    "os"
    )

func main() {
    var fruits uint
    fmt.Print("Amount of fruits: ")
    if _, err := fmt.Scan(&fruits); err!= nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    
    fmt.Println("Apple Pies:", 
        makeApplePies(fruits))
}


const (
    APPLES_PER_BOWL uint = 2
    APPLES_PER_PIE  uint = 3
)

func makeApplePies(fruitAmount uint) (pies uint) {
    apples := fruitAmount / APPLES_PER_BOWL
    pies   = apples / APPLES_PER_PIE
    return
}
