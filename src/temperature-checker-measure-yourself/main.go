/*
Measure Yourself
Temperature Checker +10 XP

You are writing a program for a temperature checking system at an airport.
The system measures the body temperature of a person and needs to output "Allowed" if it is in the normal range, or "Fever" if it is higher than normal.
Normal: up to 99.5 °F 
Fever: > 99.5 °F

The program should take the temperature as a float from input and output the corresponding message.
Sample Input:
101.3

Sample Output:
Fever

Hint:
Use an if/else statement to make the decision.


@author davorpatech
@since 2021-05-12
*/

package main

import . "fmt"

const (
    feverThreshold = 99.5
    sAllowed = "Allowed"
    sFever = "Fever"
)

func main() {
    //your code goes here
    var temperature float64
    if _, err := Scanln(&temperature); err != nil {
        panic(err)
    } else if temperature < feverThreshold {
        Print(sAllowed)
    } else {
        Print(sFever)
    }
}
