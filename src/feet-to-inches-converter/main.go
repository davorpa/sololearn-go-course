/*
Feet to Inches Converter +10 XP

You need to create a program that converts feet to inches and outputs the resulting value.
The feet value is provided as the input to the program.
1 foot = 12 inches. 

Sample Input:
5

Sample Output:
60

Hint:
The input is an integer.

@author davorpatech
@since 2021-05-12
*/
package main

import . "fmt"

const FootInches = 12

func main(){
    //your code goes here
    var feet int
    if _, err := Scanf("%d", &feet); err != nil {
        panic(err)
    }
    inches := feet * FootInches
    Println(inches)
}
