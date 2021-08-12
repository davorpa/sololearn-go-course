/**
Fibonacci series without recursion.

Using bitwise operators and a fast raising matrix algorithm.

Algorithm Complexity:
    This :  O(log n)  💯💯🔝
Recursive:  O(2^n)    😲🚫

@author davorpatech
@since  2021-07-09
*/
package main

import ."fmt"

func main() {
    for n := -500; n <= 500; n++ {
        Printf("Fib(%4d)= %30d\n", n, Fibonacci(n))
    }
}


const (
    ZERO int = 0
    ONE  int = 1
    TWO  int = 2
)

// Fibonacci(index) returns Fibonacci sequence member with corresponding index:
// - Fibonacci(0) == 0
// - Fibonacci(1) == Fibonacci(2) == 1
// - Fibonacci(3) == 2
// - Fibonacci(4) == 3
// ...
// Negative indexes produce results, extended for negative values:
// - Fibonacci(-1) == 1
// - Fibonacci(-2) == -1
// - Fibonacci(-3) == 2
// - Fibonacci(-4) == -3
// ...
// Result is calculated via matrix ((1, 1), (1, 0)) fast raising to `index` power.
func Fibonacci(index int) int {
    v0, v1 := ZERO, ONE // Result vector
    if index < ZERO {
        index = -index
        v1 = -ONE       // save sign
    }

    // init matrix ((1, 1), (1, 0))
    for m00, m01, m10, m11 := v1, v1, v1, v0;
            index != ZERO; {
        if IsOdd(index) {
            // If power is odd,
            // then multiply result vector by matrix
            v0, v1 = multMfib(v0, m00, v1, m10),
                     multMfib(v0, m01, v1, m11)
        }
        // Square the matrix
        m00, m01, m10, m11 = multMfib(m00, m00, m01, m10),
                             multMfib(m00, m01, m01, m11),
                             multMfib(m10, m00, m11, m10),                             
                             multMfib(m10, m01, m11, m11)
        // `index` fast division by 2 (moving bits)
        index >>= ONE
//        index /= TWO
    }
    return v0
}


// Check if a number is odd.
func IsOdd(value int) bool {
    // `value` fast div remainder by 2 (masking bits)
    return (value & ONE) != ZERO
//    return (value % TWO) != ZERO
}


func multMfib(v1, m1, v2, m2 int) int {
    return (v1 * m1) + (v2 * m2)
}
