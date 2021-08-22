/*
TEXTBOOKS FOR STUDENTS (+10 XP)
===============================

We need to distribute math and history textbooks to students.
There are 2 class sections: the first section has 18 pupils, and the second one has 19. The total number of books available for distribution is 76.

Write a program to calculate and output how many books will be left after each student receives both books.

TIP:
Use the modulo operator (%) to get the remainder.


@author davorpatech
@since  2021-08-23
*/

package main

import ."fmt"

func main() {
    var (
        AVAILABLE_BOOKS uint = 76
        TEXTBOOKS_PER_STUDENT uint = 2
        CLASSROOMS = []uint{18, 19}
    )
    
    distributedBooksPerClass := ApplyFactorUints(
        CLASSROOMS, 
        TEXTBOOKS_PER_STUDENT)
    
    remainingBooksAmount := AVAILABLE_BOOKS % SumUints(distributedBooksPerClass...)

    Printf(`INPUT:
  Books Amount:       %d
  Books Per Students: %d
  Classrooms:         %v
OUTPUT:
  Remaining books:    %d
`,
        AVAILABLE_BOOKS,
        TEXTBOOKS_PER_STUDENT,
        CLASSROOMS,
        remainingBooksAmount)
}


func ApplyFactorUints(values []uint, factor uint) (result []uint) {
    for _, value := range values {
        result = append(result, value * factor)
    }
    return
}

func SumUints(values... uint) (sum uint) {
    for _, value := range values {
        sum += value
    }
    return
}
