/*
Implements a Permutation Generator pattern using High Order Functions (HOF).

@author davorpatech
@since  2021-09-21
*/
package main

import "fmt"

type IntIteratorFunc func() (int, bool)
type IntPermuter func(int) int

// This function `IntPermutator` returns another function, which we define anonymously in the body of `permutator`.
// The returned `iterator` function _closes over_ the variable `data` to form a closure.
func IntPermutator(data int, permutation IntPermuter, bound int) (it IntIteratorFunc) {
    it = func() (int, bool) {
        hasNext := data < bound
        if hasNext {
            data = permutation(data)
        }
        hasNext = data < bound
        return data, hasNext
    }
    return
}

//
// permutation factory providers
//

func Incrementor(val int) (p IntPermuter) {
    p = func(j int) int {
        j += val
        return j
    }
    return
}

func Multiplicor(val int) (p IntPermuter) {
    p = func(j int) int {
        j *= val
        return j
    }
    return
}

func main() {
    var next IntIteratorFunc

    // We call `IntPermutator`, assigning the result (a function) to `next`.
    // This function value captures its own `data` value, which will be updated each time we call `next`.
    next = IntPermutator(1, Multiplicor(2), 7)
    // See the effect of the closure by calling `next` a few times.
    fmt.Println(next())
    fmt.Println(next())
    fmt.Println(next())
    fmt.Println(next())

    fmt.Println()
    
    // To confirm that the state is unique to that particular function, create and test a new ones.

    next = IntPermutator(11, Incrementor(2), 17)
    // decline first value (while loop)
    for value, hasNext := next(); hasNext; value, hasNext = next() {
        fmt.Println(value)
    }

    fmt.Println()

    next = IntPermutator(11, Incrementor(2), 17)
    // accept first value (do-while loop)
    for value, hasNext := 0, true; hasNext; {
        value, hasNext = next()
        fmt.Println(value)
    }
}
