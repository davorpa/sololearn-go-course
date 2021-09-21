/*
Flexible Fibonacci generator composited with closures, channels and High Order Functions (HOF)

@author davorpatech
@since  2021-09-21
*/
package main

import (
    "context"
    "fmt"
    "math/big"
    "time"
)

// `Fibonacci` returns a closure that returns successive Fibonacci's serie numbers.
func Fibonacci() func() *big.Int {
    a, b := big.NewInt(0), big.NewInt(1)
    return func() *big.Int {
        a, b = new(big.Int).Set(b),
               new(big.Int).Add(a, b)
        // local cloned
        return new(big.Int).Set(a)
    }
}

// `FibonacciGen` returns a channel where send the generated serie values while complains `takeWhile` predicate.
func FibonacciGen(takeWhile BintPredicate) (chan *big.Int) {
    if takeWhile == nil {
        panic("Predicate `takeWhile` not provided")
    }
    ch, fib := make(chan *big.Int), Fibonacci()
    go func() {
        defer close(ch)
        for i, x := uint(0), fib(); takeWhile(x, i); {
            // emit
            ch <- x
            // next
            i, x = i + 1, fib()
        }
    }()
    return ch
}

// `BintPredicate` defines the HOF predicate while iterate `big.Int` generators.
type BintPredicate func(*big.Int, uint) bool

// `UntilCancel` defines a predicate that returns `false` when context is cancelled.
func UntilCancel(ctx context.Context) (p BintPredicate) {
    if ctx == nil {
        panic("Cancel `ctx` not provided")
    }
    state := true
    p = func(v *big.Int, i uint) bool {
        return state
    }
    // monitor cancelling to change state
    go func(state *bool) {
        select {
            case <- ctx.Done():
                *state = false
                return
        }
    }(&state)
    return
}

// `FirstN` defines a predicate that returns `false` when n-items are reached.
func FirstN(n uint) (p BintPredicate) {
    p = func(v *big.Int, i uint) bool {
        return i < n
    }
    return
}

// `LessThanN` defines a predicate that returns `false` when a value is reached.
func LessThanN(n int) (p BintPredicate) {
    bn := big.NewInt(int64(n))
    return LessThan(bn)
}

// `LessThan` defines a predicate that returns `false` when a value is reached.
func LessThan(n *big.Int) (p BintPredicate) {
    if n == nil {
        panic("Upper bound number `n` not provided")
    }
    p = func(v *big.Int, i uint) bool {
        return v.Cmp(n) < 0
    }
    return
}

//
// Main program
//
func main() {
    fib := Fibonacci()
    // Function calls are evaluated left-to-right.
    fmt.Println(
        fib(), fib(), fib(), fib(), fib())

    fibprint(
        "First 10 Fibonacci numbers",
        FirstN(10))

    fibprint(
        "Fibonacci numbers less than 500",
        LessThanN(500))

    ctx, _ := context.WithTimeout(
        context.Background(),
        10 * time.Millisecond)
    fibprint(
        "Fibonacci during 10ms",
        UntilCancel(ctx))
}

//
// print helper
//
func fibprint(name string, takeWhile BintPredicate) {
    fmt.Println()
    fmt.Println(name)
    count := uint(0)
    for n := range FibonacciGen(takeWhile) {
        fmt.Println(n)
        count++
    }
    fmt.Println("==>", count, "numbers.")
}
