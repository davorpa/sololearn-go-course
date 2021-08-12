/*
Atomic concurrent processing

As you can see, the result is unpredictable and also there are some side effects on chained operations between goroutines, but each of those operations are guaranteed to be atomic.

Atomicity are implemented here using semaphore channels.

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
    n.Print() // 5
    
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



// A SemaphoricInt uses a binary channel to protect its  "int" wrapped value from concurrent access made over any of its operations.
type semaphoricInt struct {
    // a binary channel guarding value
    semaphore    chan struct{}
    value        int
}

// Constructor
func AtomicInt(value int) *semaphoricInt {
    n := new(semaphoricInt)
    // Buffered channels block when are full.
    // So, it's important use one with size=1 to get the locking behaviour.
    n.semaphore = make(chan struct{}, 1)
    n.value = value
    return n
}

// Implements fmt.GoStringer
func (n *semaphoricInt) GoString() string {
    // hide semaphore
    return fmt.Sprintf("%T{value:%v}", n, n.Get())
}

// Implements fmt.Stringer
func (n *semaphoricInt) String() string {
    // hide semaphore
    return fmt.Sprintf("&{%v}", n.Get())
}

func (n *semaphoricInt) Print() *semaphoricInt {
    return n.PrintIn(os.Stdout)
}

func (n *semaphoricInt) PrintIn(w io.Writer) *semaphoricInt {
    if (w == nil) {
        panic("printer: nil writer")
    }
    fmt.Fprintln(w, n.Get())
    return n
}

func (n *semaphoricInt) Get() int {
    n.semaphore <- struct{}{} // acquire token
    value := n.value
    <- n.semaphore // release token
    return value
}

func (n *semaphoricInt) Set(value int) *semaphoricInt {
    n.semaphore <- struct{}{} // acquire token
    log.Println("SET", n.value, "to", value)
    n.value = value
    <- n.semaphore // release token
    return n
}

func (n *semaphoricInt) Add(delta int) *semaphoricInt {
    n.semaphore <- struct{}{} // acquire token
    value := n.value + delta
    log.Println("ADD", delta, "to", n.value, "=", value)
    n.value = value
    <- n.semaphore // release token
    return n
}

func (n *semaphoricInt) Increment() *semaphoricInt {
    n.semaphore <- struct{}{} // acquire token
    value := n.value + 1
    log.Println("INC", 1, "to", n.value, "=", value)
    n.value = value
    <- n.semaphore // release token
    return n
}

func (n *semaphoricInt) Decrement() *semaphoricInt {
    n.semaphore <- struct{}{} // acquire token
    value := n.value - 1
    log.Println("DEC", 1, "to", n.value, "=", value)
    n.value = value
    <- n.semaphore // release token
    return n
}
