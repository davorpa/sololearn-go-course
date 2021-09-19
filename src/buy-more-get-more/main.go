/*
Buy More, Get More
else if +10 XP

A store is offering a tiered discounts based on the purchase total.

5000 and above => 50%
3000-4999      => 30%
1000-2999      => 10%
1-999          => 0%

Write a program that takes the total price as input and outputs the corresponding discount to the console.

Sample Input
4700

Sample Output
30%

Hint:
Use logical operator && to chain multiple conditions.


@author davorpatech
@since  2020-09-18
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var totalPrice float64
    fmt.Print("Enter purchase cost $")
    if _, err := fmt.Scanln(&totalPrice); err != nil || totalPrice < 0 {
        fmt.Fprintln(os.Stderr, "ERR!", "Invalid number.")
        os.Exit(1)
    }
    fmt.Println()

    fmt.Printf("The purchase discount for $%v is: %v%%\n",
        totalPrice,
        ResolvePurchaseDiscount_MAP(totalPrice))
}

func ResolvePurchaseDiscount_MAP(totalPrice float64) (discount uint) {
    if totalPrice < 0 {
        // avoid below overflow
        return
    }
    //              0   1   2   3   4  5+
    bands := []uint{0, 10, 10, 30, 30, 50}
    // apply factor to adapt/map to bands
    position := int(totalPrice / 1000)
    if l := len(bands); position >= l {
        // avoid above overflow (5+)
        position = l - 1
    }
    discount = bands[position]
    return
}

func ResolvePurchaseDiscount_IFER(totalPrice float64) uint {
    /* Sorted-ifs with early returns */
    if totalPrice >= 5000 {
        return 50
    }
    if totalPrice >= 3000 {
        return 30
    }
    if totalPrice >= 1000 {
        return 10
    }
    return 0
}

func ResolvePurchaseDiscount_IF(totalPrice float64) (discount uint) {
    /* Sorted-ifs avoid repeat pointcuts */
    if totalPrice >= 5000 {
        discount = 50
    } else if totalPrice >= 3000 {
        discount = 30
    } else if totalPrice >= 1000 {
        discount = 10
    }
    return
}

func ResolvePurchaseDiscount_ELIF(totalPrice uint) (discount uint) {
    if totalPrice >= 1 && totalPrice <= 999 {
        discount = 0
    } else if totalPrice >= 1000 && totalPrice <= 2999 {
        discount = 10
    } else if totalPrice >= 3000 && totalPrice <= 4999 {
        discount = 30
    } else if totalPrice >= 5000 {
        discount = 50
    }
    return
}
