/**
Fibonacci series without recursion.

Using bitwise operators and a fast raising matrix algorithm.

Algorithm Complexity:
    This :  O(log n)  💯💯🔝
Recursive:  O(2^n)    😲🚫

Big-Endian version (no digit limit / number overflow)

@author davorpatech
@since  2021-07-11
*/
package main

import ."fmt"
import "math/big"

func main() {
    for n, max := big.NewInt(-500), big.NewInt(500); n.Cmp(max) <= 0; n.Add(n, ONE) {
        Printf("Fib(%4d)= %30d\n", n, Fibonacci(n))
    }
}


var (
    ZERO = new(big.Int)
    ONE  = big.NewInt(1)
    TWO  = big.NewInt(2)
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
func Fibonacci(index *big.Int) *big.Int {
    if index == nil {
        return index
    }
    // Result vector (0, 1)
    v0, v1 := CloneBigInt(ZERO), CloneBigInt(ONE)
    index = CloneBigInt(index) // clone local
    if index.Cmp(ZERO) < 0 {
        index.Neg(index)
        v1.Neg(v1)        // -1. save sign
    }

    // init matrix ((1, 1), (1, 0))
    for m00, m01, m10, m11 := v1, v1, v1, v0;
            index.Cmp(ZERO) != 0; {
        if IsOdd(index) {
            // If power is odd,
            // then multiply result vector by matrix
            v0, v1 = multMfib(v0, m00, v1, m10),
                     multMfib(v0, m01, v1, m11)
        }
        // Square the matrix
        m00, m01, m10, m11 =
            multMfib(m00, m00, m01, m10),
            multMfib(m00, m01, m01, m11),
            multMfib(m10, m00, m11, m10),                             
            multMfib(m10, m01, m11, m11)
        // `index` fast division by 2 (moving bits)
        index.Rsh(index, 1)   // index >> 1
//        index.Div(index, TWO) // index /= 2        
    }
    return v0
}


// Check if a number is odd.
func IsOdd(value *big.Int) bool {
    r := new(big.Int)
    // `value` fast div remainder by 2 (masking bits)
    return r.And(value, ONE).Cmp(ZERO) != 0 // value & 1
//    return r.Mod(value, TWO).Cmp(ZERO) != 0 // value % 2
//    return r.Rem(value, TWO).Cmp(ZERO) != 0 // value % 2
}


func multMfib(v1, m1, v2, m2 *big.Int) *big.Int {
    // (v1 * m1) + (v2 * m2)
    r1, r2 := new(big.Int).Mul(v1, m1), new(big.Int).Mul(v2, m2)
    return r1.Add(r1, r2)
}

func CloneBigInt(value *big.Int) *big.Int {
    if value != nil {
        value = new(big.Int).Set(value) // clone local
    }
    return value
}
