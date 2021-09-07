/**
Big Sale. Math Operators +10 XP

Time to go shopping!
Everything in the store has been discounted by 20%.
You are given a program that takes the price of an item as input. Complete the program so that it outputs the discounted price to the console.

Sample Input
100

Sample Output
80

Explanation
20 percent of 100 equals to 20 (100 * 20/100), so the discounted price will be 80 (100 - 20).

Hint
Remember the division (/) and multiplication (*) operators.


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
    fmt.Print("Price? ")
    if _, err := fmt.Scanln(&text); err != nil {
        fmt.Fprintln(os.Stderr, "ERR", err)
        os.Exit(1)
    }
    
    price, err := strconv.ParseFloat(text, 64)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR", err)
        os.Exit(1)
    }
    
    price = applyDiscount(price)
    fmt.Printf("\nDiscounted: %v\n", price)
}

func calculateDiscount(price float64) float64 {
    return price * 0.2
}

func applyDiscount(price float64) float64 {
    discount := calculateDiscount(price)
    return price - discount
}
