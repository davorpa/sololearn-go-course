/*
Spin Lock concurrent pattern implementation applied to ensure atomicity of "int" operations.

As you can see, the result is unpredictable and also there are some side effects on chained operations between goroutines, but each of those operations are guaranteed to be atomic.


This lock type implemented using internally the CAS algorithm (compare and swap)

CAS algorithm is a well-known unlocked algorithm. Lock-free programming is to synchronize variables between threads without using locks, that is, to synchronize variables without thread blocking, so it is also called non-blocking Synchronization. The CAS algorithm involves three operands

a) Memory value V to read and write
b) Value A for comparison
c) New value B to be written

If and only if the value of V equals A, CAS updates the value of V with the new value B atomically, otherwise no operation will be performed (comparison and substitution is an atomic operation). In general, it is a spin operation, i.e. repeated retries.


@author davorpatech
@since 2021-06-14
*/

package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "sync"
    "sync/atomic"
    "runtime"
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



/*
Spinlock refers to when a thread acquires a lock, if the lock has been acquired by other threads, then the thread will wait for a loop, and then constantly determine whether it can be successfully acquired, until the lock is acquired, it will not exit the loop.

Threads that acquire locks are always active, but do not perform any valid tasks. Using such locks results in busy-waiting.

It is a lock mechanism proposed to protect shared resources. In fact, spin locks are similar to mutexes in that they are used to solve the mutex of a resource. Whether mutex or spin lock, at any time, only one holder can get the lock at most, that is to say, only one execution unit can get the lock at any time. But they are slightly different in scheduling mechanism. For mutexes, if resources have been occupied, resource applicants can only sleep. But the spin lock does not cause the caller to sleep. If the spin lock has been maintained by other executing units, the caller keeps looping there to see if the holder of the spin lock has released the lock. That is why the term “spin” is named.


Then...When to SpinLock? When to Mutex?
Differences?

* SpinLock: As you can see, SpinLock idea comes from Loop-CAS, which is lock-free, but kinda busy checking.
* Mutex: When Lock() invocation happens, however can’t get the lock, the thread/goroutine will be added to wait list of scheduler and thread/goroutine switch happens at this time.

So SpinLock thinks that I should very probably get the lock in my scheduling CPU slot, and in a aggressive way.

What about Mutex scheduling? In Golang implementation, it’s:

1. If Lock() happens and can’t get the lock, will be added to wait list,
2. For new young goroutine who already get the lock, will just let them run, since they already get the lock,
3. If the old goroutine has been wait for very long time, will give the lock to him and let him run, so it won’t wait for too long. By default, starvationThresholdNs is 1e6 nanoseconds (1ms). Normally it’s FIFO order, but if starvationThresholdNs meet, will enter starvation mode.

In a summary:

1. Under multi-core CPU, if you think that you will quickly get the lock, then please use SpinLock, so that it doesn’t need context switch to get the lock and execute the code.
  1a. the time to get the lock is short (short task or low concurrency)
  1b. the critical code block running time cost won’t be too much, i.e., can quickly do the shared variables computation and release the lock
  1c. the concurrency is not too much, or it will needs lots of context switches to release the lock

2. Under single-core CPU, never use SpinLock, since it will try to use up the CPU slot, until context switch happens and then back to get the lock. Context Switch always happens in 1 Core CPU. But the above SpinLock code implementation is good, since it always yield the CPU to the other goroutines by using runtime.Gosched(). So the SpinLock implementation here is a bit like a Mutex but lock-free.

3. Under multi-core CPU, if you think your critical code block can’t complete in the assigned CPU slot, please try to use Mutex, so that it won’t waste up much time on getting the lock.
  3a. the time to get the lock is long (long task or high concurrency)
  3b. the critical code block running time cost is little large, i.e, keeps computation on top of shared variables for a long time. So it’s important not to put I/O into critical code block.

To sum up:

1. Single-core: not use SpinLock, or at least yield the processor every time failed CAS.
2. Multi-core, short task or low concurrency: use SpinLock.
3. Multi-core, long task or high concurrency: use Mutex.
4. Don't use I/O in critical code block, this should be a common knowledge.

If high concurrency with heavy I/O, then better to use goroutine + channel solution.
*/
type spinLock uint32

func SpinLock() sync.Locker {
    return new(spinLock)
}

// Lock adquires the spin lock and then blocks other goroutines until it is released.
func (sl *spinLock) Lock() {
    // keep checking, until sl == 0, i.e, nobody locks it,
    // and then change the sl = 1, i.e, lock it
    for !sl.TryLock() {
        // yield goroutine while it cannot be adquired
        runtime.Gosched()
        // Gosched yields the processor, allowing other goroutines to run.
        // It does not suspend the current goroutine, so execution resumes automatically.
    }
}

// Unlock releases the spin lock adquired in some goroutine, giving the chance to others adquire it.
// Unlike [Mutex.Unlock](http://golang.org/pkg/sync/#Mutex.Unlock), there's no harm calling it on an unlocked SpinLock.
func (sl *spinLock) Unlock() {
    // set atomically to zero releasing the lock and then give the opportunity to others goroutines adquire it (breaking its CAS loop)
    atomic.StoreUint32((*uint32)(sl), 0)
}

// TryLock will try to lock sl and return whether it succeed or not without blocking.
func (sl *spinLock) TryLock() bool {
    // change state atomically from 0 to 1
    return atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1)
}

func (sl *spinLock) String() string {
    if atomic.LoadUint32((*uint32)(sl)) == 1 {
        return "Locked"
    }
    return "Unlocked"
}






// A LockedInt uses any "sync.Locker" to protect its  "int" wrapped value from concurrent access made over any of its operations.
type lockedInt struct {
    lk        sync.Locker
    value     int
}

// Constructor
func AtomicInt(value int) *lockedInt {
    n := new(lockedInt)
    n.value = value
    n.lk = SpinLock()
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
    n.lk.Lock() // acquire token
    defer n.lk.Unlock() // release token
    return n.value
}

func (n *lockedInt) Set(value int) *lockedInt {
    n.lk.Lock() // acquire token
    defer n.lk.Unlock() // release token
    log.Println("SET", n.value, "to", value)
    n.value = value 
    return n
}

func (n *lockedInt) Add(delta int) *lockedInt {
    n.lk.Lock() // acquire token
    defer n.lk.Unlock() // release token
    value := n.value + delta
    log.Println("ADD", delta, "to", n.value, "=", value)
    n.value = value
    return n
}

func (n *lockedInt) Increment() *lockedInt {
    n.lk.Lock() // acquire token
    defer n.lk.Unlock() // release token
    value := n.value + 1
    log.Println("INC", 1, "to", n.value, "=", value)
    n.value = value
    return n
}

func (n *lockedInt) Decrement() *lockedInt {
    n.lk.Lock() // acquire token
    defer n.lk.Unlock() // release token
    value := n.value - 1
    log.Println("DEC", 1, "to", n.value, "=", value)
    n.value = value
    return n
}
