/*
8s Alive program. Let's hack Sololearn 1s limit.

TOPIC: Signals & I/O flood

Signals received by SoloLearn Docker instance provided by Azure container:
- 15 SIGTERM: kill attempt after 1 second
- 21 SIGURG:  too much resources used

@author davorpatech
@since 2021-05-16
*/
package main

import (
    "log"
    "math/rand"
    "os"
    "os/signal"
    "syscall"
    "time"
)

//======================================
// PROGRAM
//

var (
    duration      time.Duration
    floodInterval time.Duration
    gcgenInterval time.Duration
)


func init() { // init this package
    // SL playground fight
    duration      = 7700 * time.Millisecond
    floodInterval =  100 * time.Millisecond
    gcgenInterval =  200 * time.Millisecond
    
    // configure default logger
    log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC | log.Lmsgprefix)
    log.SetOutput(os.Stdout)
    
    // configure default randomizer to be non-constant
    rand.Seed(time.Now().UnixNano())
}


func main() {
    ender := traceStartEnd(time.Now())
    defer ender()
    
    // Set up the waiter channel
    waiter := time.NewTimer(duration)
    // ensure release resources on exit
    defer waiter.Stop()
    
    // Setup the flood ticker
    fluzer := time.NewTicker(floodInterval)
    // ensure release resources on exit
    defer fluzer.Stop()
    
    // Setup the GC entropy ticker
    gcgenr := time.NewTicker(gcgenInterval)
    defer gcgenr.Stop()
    
    // Set up channel on which to send signal notifications.
    // We must use a buffered channel or risk missing the signal if we're not ready to receive when the signal is sent.
    sigChan := make(chan os.Signal, 1)
    // ensure stop listening on exit
    defer signal.Stop(sigChan)
    // Passing no signals to Notify means that all signals will be sent to the channel.
    signal.Notify(sigChan)
    // Ignore some signals (if no values are passed,  means all incoming signals will be ignored on Notify).
    if sig := syscall.SIGURG; rand.Intn(1000) > 500 {
        log.Printf("☑ Ignoring signal: %[1]d-%[1]v", sig)
        // ignore non-filtered Docker-Go! signal
        signal.Ignore(sig)
    } else {
        log.Printf("☒ Signal \"%[1]d-%[1]v\" not filtered.", sig)
    }
    
    // channels listening
    LOOP:
    for numFluz,numSig := 0,0; ; {
        // receiver selector
        select {
            // catch program end
            case <- waiter.C:
                break LOOP
            // catch process signals
            case s := <- sigChan:
                numSig += 1
                var ico string
                switch s {
                    case syscall.SIGURG: ico = "⚠"
                    case syscall.SIGTERM: ico = "🏁⚠"
                    default: ico = "🐙"
                }
                log.Printf("%v Sig %d: %d-%v\n", ico, numSig, s, s)
            // catch flood ticker
            case t := <- fluzer.C:
                _ = t
                numFluz += 1
                log.Println("😈 Flood", numFluz, "...")
            // generate some GC activities to produce more 21 SIGURG signals
            case <- gcgenr.C:
                _ = new(int)
        }
    }
}


func traceStartEnd(start time.Time) (ender func()) {
    log.Println("ℹ🚩 Start...")
    return func() {
        log.Println("ℹ🏁 End...", time.Since(start))
    }
}
