/*
Atomic concurrent processing

As you can see, the result is unpredictable and also there are some side effects on chained operations between goroutines, but each of those operations are guaranteed to be atomic.

Atomicity are implemented here using "sync.Mutex"

@author davorpatech
@since 2021-06-10
*/

package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "sync"
)

func init() {
    log.SetFlags(log.Ltime|log.Lmicroseconds|log.LUTC)
    log.SetOutput(os.Stdout)
}

func main() {
    n := AtomicInt(5)
    n.Print()  // 5
    
    var wg sync.WaitGroup
    wg.Add(4) // the goroutines amount
    
    go func() {
        defer wg.Done()
        n.Decrement().Decrement().Print()
    }()
    go func() {
        defer wg.Done()
        n.Set(15).Print()
    }()
    go func() {
        defer wg.Done()
        n.Add(10).Increment().Print()
    }()
    go func() {
        defer wg.Done()
        n.Add(-8).Print()
    }()
    
    // wait for all goroutines are done
    wg.Wait()
    fmt.Println("END", n.Get())
    fmt.Println("fmt.Stringer:", n)
    fmt.Printf("fmt.GoStringer: %#v", n)
}



// A LockedInt uses a "sync.Mutex" to protect its  "int" wrapped value from concurrent access made over any of its operations.
type lockedInt struct {
    mu        sync.Mutex
    value     int
}

// Constructor
func AtomicInt(value int) *lockedInt {
    n := new(lockedInt)
    n.value = value
    return n
}

// Implements fmt.GoStringer (%#v verb)
func (n *lockedInt) GoString() string {
    // hide mutex
    return fmt.Sprintf("%T{value:%v}", n, n.Get())
}

// Implements fmt.Stringer
func (n *lockedInt) String() string {
    // hide mutex
    return fmt.Sprintf("&{%v}", n.Get())
}

func (n *lockedInt) Print() *lockedInt {
    return n.PrintIn(os.Stdout)
}

func (n *lockedInt) PrintIn(w io.Writer) *lockedInt {
    if (w == nil) {
        panic("printer: nil writer")
    }
    fmt.Fprintln(w, n.Get())
    return n
}

func (n *lockedInt) Get() int {
    n.mu.Lock() // acquire token
    defer n.mu.Unlock() // release token
    return n.value
}

func (n *lockedInt) Set(value int) *lockedInt {
    n.mu.Lock() // acquire token
    defer n.mu.Unlock() // release token
    log.Println("SET", n.value, "to", value)
    n.value = value 
    return n
}

func (n *lockedInt) Add(delta int) *lockedInt {
    n.mu.Lock() // acquire token
    defer n.mu.Unlock() // release token
    value := n.value + delta
    log.Println("ADD", delta, "to", n.value, "=", value)
    n.value = value
    return n
}

func (n *lockedInt) Increment() *lockedInt {
    n.mu.Lock() // acquire token
    defer n.mu.Unlock() // release token
    value := n.value + 1
    log.Println("INC", 1, "to", n.value, "=", value)
    n.value = value
    return n
}

func (n *lockedInt) Decrement() *lockedInt {
    n.mu.Lock() // acquire token
    defer n.mu.Unlock() // release token
    value := n.value - 1
    log.Println("DEC", 1, "to", n.value, "=", value)
    n.value = value
    return n
}
