/*
Atomic concurrent processing

As you can see, the result is unpredictable and also there are some side effects on chained operations between goroutines, but each of those operations are guaranteed to be atomic.

Atomicity are implemented here using the Compare-And-Swap (CAS) loop pattern until right value is loaded/stored into the internal "unsafe.Pointer" value.

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
    "sync/atomic"
    "unsafe"
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
    
    go func(done func()) {
        defer done()
        n.Decrement().Decrement().Print()
    }(wg.Done)
    go func(done func()) {
        defer done()
        n.Set(15).Print()
    }(wg.Done)
    go func(done func()) {
        defer done()
        n.Add(10).Increment().Print()
    }(wg.Done)
    go func(done func()) {
        defer done()
        n.Add(-8).Print()
    }(wg.Done)
    
    // wait for all goroutines are done
    wg.Wait()
    fmt.Println("END", n.Get())
    fmt.Println("fmt.Stringer:", n)
    fmt.Printf("fmt.GoStringer: %#v", n)
}



// An AtomicInt uses a "unsafe.Pointer" as internal "int" value holder and guards from concurrent access made over any of its operations using the low level facilities provided by "sync/atomic" module.
type atomicInt struct {
    value     unsafe.Pointer  // *int
}

// Constructor
func AtomicInt(value int) *atomicInt {
    n := new(atomicInt)
    p := unsafe.Pointer(&value)
    atomic.StorePointer(&n.value, p)
    return n
}

// Implements fmt.GoStringer (%#v verb)
func (n *atomicInt) GoString() string {
    // resolve value, hiding unsafe pointer
    return fmt.Sprintf("%T{value:%v}", n, n.Get())
}

// Implements fmt.Stringer
func (n *atomicInt) String() string {
    // resolve value, hiding unsafe pointer
    return fmt.Sprintf("&{%v}", n.Get())
}

func (n *atomicInt) Print() *atomicInt {
    return n.PrintIn(os.Stdout)
}

func (n *atomicInt) PrintIn(w io.Writer) *atomicInt {
    if (w == nil) {
        panic("printer: nil writer")
    }
    fmt.Fprintln(w, n.Get())
    return n
}

func (n *atomicInt) Get() int {
    // 1. load unsafe pointer value by it address
    // 2. unsafe convert to int pointer
    // 3. get value that int pointer points to
    return *(*int)(atomic.LoadPointer(&n.value))
}

func (n *atomicInt) Set(value int) *atomicInt {
/*
    vp := unsafe.Pointer(&value)
    // ensure atomicity using a Compare-And-Swap loop 
    for {
        np := atomic.LoadPointer(&n.value)
        old := *(*int)(np)
        if atomic.CompareAndSwapPointer(&n.value, np, vp) {
            log.Println("SET", old, "to", value)
            break
        }
    }
*/
    casLoopPointerInt(&n.value,
        func(old int) int { return value },
        func(old, new int) {
            log.Println("SET", old, "to", new)
        })
// or weak set with...
//    atomic.StorePointer(&n.value, unsafe.Pointer(&value))
    return n
}

func (n *atomicInt) Add(delta int) *atomicInt {
    casLoopPointerInt(&n.value,
        func(old int) int { return old + delta },
        func(old, new int) {
            log.Println("ADD", delta, "to", old, "=", new)
        })
    return n
}

func (n *atomicInt) Increment() *atomicInt {
    casLoopPointerInt(&n.value,
        func(old int) int { return old + 1 },
        func(old, new int) {
            log.Println("INC", 1, "to", old, "=", new)
        })
    return n
}

func (n *atomicInt) Decrement() *atomicInt {
    casLoopPointerInt(&n.value,
        func(old int) int { return old - 1 },
        func(old, new int) {
            log.Println("DEC", 1, "to", old, "=", new)
        })
    return n
}

func casLoopPointerInt(
        addr *unsafe.Pointer,
        evalFunc  func(old int) int, 
        afterFunc func(old, new int)) {
    // ensure atomicity using a Compare-And-Swap loop 
    for {
        // get unsafe pointer by their pointer address
        np := atomic.LoadPointer(addr)
        // unsafe cast to int pointer and get value
        old := *(*int)(np)
        // operate
        value := evalFunc(old)
        vp := unsafe.Pointer(&value)
        // loop until old and new values match
        if atomic.CompareAndSwapPointer(addr, np, vp) {
            if afterFunc != nil { // optional callback
                afterFunc(old, value)
            }
            return
        }
    }
}
