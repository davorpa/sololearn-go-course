/*
Let's Get Loud.
Sound Frequencies +10 XP

The human ear can hear sounds with frequencies in the range of 20 Hz to 20,000 Hz.
This is called the Audible frequency.
Anything below 20 Hz is called Infrasound, while anything above 20,000 Hz is called Ultrasound.

You need to make a program that takes a sound frequency as input and outputs the corresponding category. 

Hint:
In case the given input is negative, your program should output "Wrong Input".

Sample Input:
1800

Sample Output: 
Audible


@author davorpatech
@since 2021-05-14
*/

package main

import "fmt"

func main() {
  var f int
  fmt.Scanln(&f)
  //your code goes here
  switch {
    case f >= 20 && f <= 20000:
      fmt.Print("Audible")
    case f > 20000:
      fmt.Print("Ultrasound")
    case f < 20 && f >= 0:
      fmt.Print("Infrasound")
    default:
      fmt.Print("Wrong Input")
  }
}
