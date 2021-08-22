/*
Grasshopper - Grade book

Complete the `GetGrade` function so that it finds the average of the three scores passed to it and returns the letter value associated with that grade.

Numerical Score        Letter Grade
  90 <= score <= 100     'A'
  80 <= score < 90       'B'
  70 <= score < 80       'C'
  60 <= score < 70       'D'
  0  <= score < 60       'F'


@author davorpatech
@since  2021-08-22
@impl   no ifs/switch; best performance
*/

package main

import ."fmt"

func main() {
    Println("'should return a A'")
    Println(GetGrade(95, 90, 93) == 'A')
    Println(GetGrade(100, 85, 96) == 'A')
    Println(GetGrade(92, 93, 94) == 'A')
    Println(GetGrade(100, 100, 100) == 'A')
    Println("'should return a B'")
    Println(GetGrade(70, 70, 100) == 'B')
    Println(GetGrade(82, 85, 87) == 'B')
    Println(GetGrade(84, 79, 85) == 'B')
    Println("'should return a C'")
    Println(GetGrade(70, 70, 70) == 'C')
    Println(GetGrade(75, 70, 79) == 'C')
    Println(GetGrade(60, 82, 76) == 'C')
    Println("'should return a D'")
    Println(GetGrade(65, 70, 59) == 'D')
    Println(GetGrade(66, 62, 68) == 'D')
    Println(GetGrade(58, 62, 70) == 'D')
    Println("'should return an F'")
    Println(GetGrade(44, 55, 52) == 'F')
    Println(GetGrade(48, 55, 52) == 'F')
    Println(GetGrade(58, 59, 60) == 'F')
    Println(GetGrade(0, 0, 0) == 'F')
}




func GetGrade(a, b, c int) (grade rune) {
    score := Average(a, b, c)
    position := int(score / 10)
    grade = []rune("FFFFFFDCBAA")[position]
    return
}

func Average(values... int) (avg float64) {
    for _, v := range values {
        avg += float64(v)
    }
    avg /= float64(len(values))
    return
}
