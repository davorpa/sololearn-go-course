/*
Admission in Pairs
The if else statement +10 XP

Entrance to the club is only permitted in pairs.
Take the number of customers in the club as input, and, if all of them have a pair, output to the console "Right", otherwise output "Wrong".

Sample Input
14

Sample Output
Right

Hint:
Do not confuse the = operator with the == operator.


@author davorpatech
@since  2021-09-17
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var numberOfCustomers uint
    fmt.Print("Number of customers? ")
    if _, err := fmt.Scanln(&numberOfCustomers); err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", "Invalid number")
        os.Exit(1)
    }
    fmt.Println()

    fmt.Println(FormatPermittedEntrance(
        numberOfCustomers))
}

func IsEntrancePermitted(numberOfCustomers uint) bool {
    // in pairs
    return numberOfCustomers % 2 == 0
}

func FormatPermittedEntrance(numberOfCustomers uint) (text string) {
    if IsEntrancePermitted(numberOfCustomers) {
        text = fmt.Sprint("Right")
    } else {
        text = fmt.Sprint("Wrong")
    }
    return
}
