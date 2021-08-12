/*
↔ Pronic Numbers
https://www.sololearn.com/learn/4727/?ref=app

ℹ Pronic number is a number which is the product of two consecutive integers, that is, n=x*(x+1).
Write a program which, for a given input n, detect if the number is pronic or not.

3 implementations:
☞ Iterative.  O(n/2) linear
☞ Recursive.  O(n/2) linear
☞ No Loops.   O(1) constant

📮 EXAMPLES: 
6: True
42: True
132: True
2550: True
10100: True
4: False
66: False
133: False
1666: False
17289: False

🏳Additional challenge:
Do not use loops to perform this task!

@author davorpatech
@since  2021-07-11
*/
package main

import ."fmt"
import "math"


func main() {
    pronics   := []uint{0, 6, 42, 132, 2550, 10100}
    nopronics := []uint{1, 4, 66, 133, 1666, 17289}
    impls, kw := map[string](func(uint) bool) {
        "O(n) iterative": IsPronicIterative,
        "O(n) recursive": IsPronicRecursive,
        "O(1) direct":    IsPronicDirect,
    }, 14
    
    for name, fn := range impls {
        for _, n := range pronics {
            if fn(n) {
                Printf("✅ Pronic %*v: %5d\n", kw+2, name, n)
            } else {
                Printf("❎ Pronic %*v: %5d\n", kw+2, name, n)
            }
        }
        
        for _, n := range nopronics {
            if !fn(n) {
                Printf("✅ NoPronic %*v: %5d\n", kw, name, n)
            } else {
                Printf("❎ NoPronic %*v: %5d\n", kw, name, n)
            }
        }
        Println()
    }
}


// No loops approach with constant Big-O
// pow(x, n) inverse operation is root(x, n)
func IsPronicDirect(n uint) (pronic bool) {
    x := uint(math.Sqrt(float64(n)))
    // programatic boxing/unboxing is a must
    // why? operators doesn't work mixing types (except with number literals)
    pronic = n == x * (x + 1)
    return
}


func IsPronicIterative(n uint) (pronic bool) {
    pronic = n == 0
    // Optimized n/2. Go further away isn't needed.
    for x, max := uint(1), n/2; x < max && !pronic; x++ {
        pronic = n == x * (x+1)
    }
    return
}


func IsPronicRecursive(n uint) bool {
    return isPronicRecursiveX(n, 0)
}

func isPronicRecursiveX(n, x uint) bool {
    // base cases
    switch {
        // accept condition
        case n == x * (x+1) || n == 0: return true
        // like loop breaker
        // Optimized n/2. Go further away isn't needed.
        case x >= n/2: return false
    }
    // next case
    return isPronicRecursiveX(n, x+1)
}
