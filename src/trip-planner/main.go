/*
Trip Planner +50 XP

You need to plan a road trip. You are traveling at an average speed of 40 miles an hour.
Given a distance in miles as input (the code to take input is already present), output to the console the time it will take you to cover it in minutes.

Sample Input:
150

Sample Output:
225

Explanation:
It will take 150/40 = 3.75 hours to cover the distance, which is equivalent to 3.75*60 = 225 minutes.


@author davorpatech
@since  2021-09-16
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var distance float64
    fmt.Print("Input trip distance (miles): ")
    if _, err := fmt.Scanln(&distance); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    time := CalculateTripTime(distance, AVERAGE_SPEED)
    time = HoursToMinutes(time)

    fmt.Printf("The trip can be covered in %.5f %v\n",
        time,
        PluralizeF64(time, "minute", "minutes"))
}

func PluralizeF64(value float64, singleLiteral, pluralLiteral string) string {
    if value != 1 {
        return pluralLiteral
    }
    return singleLiteral
}

const AVERAGE_SPEED float64 = 40

func CalculateTripTime(
    distance, averageSpeed float64) float64 {
    ensurePositiveF64(distance, "distance")
    ensurePositiveF64(averageSpeed, "average speed")
    return distance / averageSpeed
}

func HoursToMinutes(t float64) float64 {
    ensurePositiveF64(t, "hours")
    return t * 60
}

func ensurePositiveF64(v float64, field string) float64 {
    if v < 0 {
        panic(fmt.Sprintf("%v must be a positive number: %v", field, v))
    }
    return v
}
