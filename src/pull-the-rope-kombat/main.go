/*
TOPIC: Extreme OOP, serial/atomic/synced counters. Secrets on Go! concurrency revealed.

This time, the problem to solve is supervise a battle between two teams pulling from a rope and determine who wins. The result is totally non deterministic.

Cancellable contexts and mutex objects or atomic operations are used to orchestrate all. The counter acts as message / internal state. Each actor (referee, players, spectators...) works in their own goroutine.

Design Patterns you can see coded here:
- Factory method + Provider
- Producer (players)
- Consumer (referee / spectators)
- Delegate / Decorator / Wrapper
- Muteable guard
- CAS (Compare-And-Swap atomic loop guard)
- Crypto Secure randoms

Clean code / SOLID:
- Interface Segregation
- OCP. Open-Close Principle
- IOC. NewXXX method signatures return ifaces
- Composition vs Inheritance
- TDD ready

@author davorpatech
@since 2021-05-27
*/

package main

import (
    "context"
    crand "crypto/rand"
    "encoding/binary"
    "encoding/json"
    _"encoding/xml"
    "fmt"
    "log"
    _"math/big"
    mrand "math/rand"
    "os"
    "os/signal"
    "runtime"
    _"strings"
    "strconv"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
    "unsafe"
)


func init() { // init this package
    // Configure default logger
    log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC | log.Lmsgprefix)
    log.SetPrefix("")
    log.SetOutput(os.Stdout)
    
    // some runtime info
    log.Print(`ℹ Version: `, runtime.Version())
    log.Print(`ℹ NumCPU: `, runtime.NumCPU())
    log.Print(`ℹ Ini GOMAXPROCS: `, runtime.GOMAXPROCS(0))
    
    // increase parallelism in goroutines
    runtime.GOMAXPROCS(runtime.NumCPU()*2) // 2x
    log.Print(`ℹ New GOMAXPROCS: `, runtime.GOMAXPROCS(0))
}



//======================================
// PROGRAM
//


func main() {
    defer StopwatchTrace("")()
    
    // Manage gracefull kill. Extends SL timeout.
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGTERM)
    defer signal.Stop(sigchan)
    
    // select shared counter by all goroutines
    meterFactory := IntCounterRandomFactory(
            NewTimeSeededRandom())
    // making an instance
    var meter IntCounter = meterFactory()
    log.Printf(`ℹ Metering with "%T"`, meter)
    
    // concurrency control
    var wg sync.WaitGroup
    bgctx := context.Background()
    
    // create actors
    player1 := NewPullRopePlayer("Team⚂",
        NewPullRopeFuncAction("⤴",
            meter.IncrementAndGet),
        NewDeltaTicker(5, 250 + uint(mrand.Intn(100))))
    player2 := NewPullRopePlayer("Team⚃",
        NewPullRopeFuncAction("⤵",
            meter.DecrementAndGet),
        NewDeltaTicker(5, 250 + uint(mrand.Intn(100))))
    
    // Start referee supervising ===========
    supervisedCtx := referee(bgctx, &wg, meter, NewFixedTicker(5))
    // Start battle metering ================
    // The callers (supervisor) need to cancel the context once
    // they are done consuming not to leak
    // the internal goroutines started inside.
    player1.Start(supervisedCtx, &wg)
    player2.Start(supervisedCtx, &wg)
    
    log.Println("💣💥 Go!! Sssshhhhhh!")
    // wait for all registered workers enď.
    wg.Wait()
    // @TODO refactor struct Game?
    printWinner(player1, player2, meter)
}


func printWinner(p1, p2 *PullRopePlayer, meter IntCounter) {
    value := meter.GetInt()
    log.Print("🎯 Metering ends: ", value)
    switch true {
        case value < 0:
            log.Print("🏆🏅 The winner: ", p2)
        case value > 0:
            log.Print("🏆🏅 The winner: ", p1)
    }
}




//======================================
// MODEL
//


// Provides the context where the battle is performed, mediating when is over and determining who wins.
//@TODO: refactor into a struct similar to Player
func referee(
        ctx context.Context, wg *sync.WaitGroup,
        meter IntCounter, 
        ticker Sleeper) context.Context {
    // create battle context
    supervisedCtx, cancel := context.WithCancel(ctx)
    wg.Add(1) // Register worker
    go func() {
        defer wg.Done() // deregister worker
        defer cancel()  // send supervisedCtx.Done()
        log.Printf("👮 Referee starts supervising every %v.", ticker)
        for {
            switch v := meter.GetInt(); {
                // incrementer wins
                case v > 5: return
                // decrementer wins
                case v < -5: return
            }
            ticker.Sleep()
        }
    }()
    return supervisedCtx
}



type PullRopePlayer struct {
    name     string
    action   NamedCaller
    ticker   Ticker
}

func NewPullRopePlayer(
        name    string,
        action  NamedCaller,
        ticker  Ticker) *PullRopePlayer {
    p := new(PullRopePlayer)
    p.name = name
    p.action = action
    p.ticker = ticker
    return p
}

func (p *PullRopePlayer) String() string {
    return fmt.Sprint(p.name)
}
func (p *PullRopePlayer) Start(
        ctx context.Context, wg *sync.WaitGroup,
        ) context.Context {
    // create battle context
    pullingCtx, cancel := context.WithCancel(ctx)
    log.Printf("ℹ %v pulls every %v.", p, p.ticker)
    wg.Add(1) // Register worker
    sleepChan, killSleeper := p.ticker.SpawnSleeper()
    go func() {
        defer wg.Done() // deregister worker
        defer cancel()  // send pullingCtx.Done()
        defer killSleeper() // stop ticker
        log.Printf("⛿ %v starts pulling!", p)
        for n := 0; ; n++ {
            select {
                // supervisor wants cancel
                case <- ctx.Done():
                    log.Printf("🏁 %v ends. %v%vpulls", p, n, p.action)
                    // returning not to leak this goroutine
                    return
                case <- sleepChan:
                    res, _ := p.action.Call()
                    log.Printf("%v pulls %v#%02v: %v", p, p.action, n, res)
            }
        }
    }()
    return pullingCtx
}



// Sleeper iface contract.
// It provides to an operation sleep capabilities a certain time interval.
type Sleeper interface {
    GetInterval() time.Duration
    Sleep()       time.Duration
}

// StartNewSleeper creates and starts a new routine that emits the sleep interval duration each time it is elapsed.
// The returned concurrent channel is where this signals are periodically notified and the function have to be called to stop / cancel the routine and then, the program does not remain in a leaked state forever.
func startNewSleeper(s Sleeper) (channel chan time.Duration, cancelFunc func()) {
    // stop control with buffered receiver channel
    done, stopping := make(chan bool, 1), false
    // unbuffered emit channel
    channel = make(chan time.Duration)
    // cancel callback
    cancelFunc = func() {
        // mark to avoid send signals to a channel in closing progress
        stopping = true
        // send stop signal
        done <- true
    }
    // start worker
    go func() {
        defer close(channel)
        for {
            select {
                // stop signal received
                case <-done:
                    // returning not to leak this goroutine
                    return
                default:
                    // wait and emit this interval
                    if stopping { return }
                    d := s.Sleep()
                    if stopping { return }
                    channel <- d
            }
        }
    }()
    return
}



// Ticker iface contract.
// A ticker allows you to sleep operations every certain time interval starting processes, as many as you want, to emit this wait interval through a concurrent communication channel.
type Ticker interface {
    GetInterval() time.Duration
    SpawnSleeper() (channel chan time.Duration, kill func())
}



// SleeperTicker iface contract.
// It joins the Sleeper and Ticker capabilities together.
type SleeperTicker interface {
    Sleeper
    Ticker
}



// A FixedTicker is a Ticker that always emits sleep intervals every certain constant milliseconds.
type FixedTicker struct {
    millis    time.Duration
}

// NewFixedTicker makes a new Ticker that always emits sleep intervals every certain constant milliseconds.
func NewFixedTicker(millis uint) SleeperTicker {
    if millis < 1 {
        panic("ticker: millis must be positive.")
    }
    t := new(FixedTicker)
    t.millis = time.Duration(millis) * time.Millisecond
    return t
}

// String implements fmt.Stringer using the fixed milliseconds interval to represent the ticker.
func (t FixedTicker) String() string {
    return fmt.Sprint(t.GetInterval())
}

// GetInterval returns the milliseconds interval which the ticker was created with.
func (t FixedTicker) GetInterval() time.Duration {
    return t.millis
}

// Sleep waits the constant milliseconds interval which the ticker was created with.
func (t FixedTicker) Sleep() (dur time.Duration) {
    dur = t.GetInterval()
    time.Sleep(dur)
    return
}

// SpawnSleeper opens a new signal emitter using the constant wait interval configured by this ticker.
// The returned concurrent channel is where this signals are periodically notified and the function have to be called to that routine stops / cancels, and then the program does not hold in a leaked state forever.
func (t FixedTicker) SpawnSleeper() (chan time.Duration, func()) {
    return startNewSleeper(&t)
}



type DeltaTicker struct {
    millis    time.Duration
    delta     uint
    r         *mrand.Rand
}

func NewDeltaTicker(millis, delta uint) SleeperTicker {
    if millis < 1 {
        panic("ticker: millis must be positive.")
    }
    if delta < 0 {
        panic("ticker: delta millis must be positive.")
    }
    t := new(DeltaTicker)
    t.millis = time.Duration(millis) * time.Millisecond
    t.delta = delta
    // Seeding with the same value results in the same random sequence each run.
    // For different numbers, seed with a different value, such as time.Now().UnixNano(), which yields a constantly-changing number.
    // For a secure seed use "crypto/rand" instead "time"
    t.r = NewSecureSeededRandom()
//    t.r = NewTimeSeededRandom()
    return t
}

// String implements fmt.Stringer using the milliseconds and delta which the ticker was created with.
func (t DeltaTicker) String() string {
    return fmt.Sprint(t.millis, " (±", t.delta, ")")
}

// GetInterval return a pseudo-random milliseconds interval calculated using a central value and a variable maximum delta which the ticker was created with.
func (t DeltaTicker) GetInterval() (dur time.Duration) {
    dur = t.millis
    // add random delta
    delta := uint(t.r.Int()) % t.delta
    dur += time.Duration(delta) * time.Millisecond
    return dur
}

func (t DeltaTicker) Sleep() (dur time.Duration) {
    dur = t.GetInterval()
    time.Sleep(dur)
    return
}

// SpawnSleeper opens a new signal emitter using the constant wait interval configured by this ticker.
// The returned concurrent channel is where this signals are periodically notified and the function have to be called to that routine stops / cancels, and then the program does not hold in a leaked state forever.
func (t DeltaTicker) SpawnSleeper() (chan time.Duration, func()) {
    return startNewSleeper(&t)
}



// A Caller is like a Runner but can return values.
type Caller interface {
    Call() (res interface{}, err error)
}

// A NamedCaller extends Caller contract to provide it a descriptible name.
type NamedCaller interface {
    Caller
    fmt.Stringer
    GetName() string
}

// Alias matching function PullRopeAction / Counter signature.
type PullRopeActionFunc func() interface{}

// PullRopeFuncAction is a NamedCaller delegating/routing the call logic strategy into a PullRopeActionFunc.
type PullRopeFuncAction struct {
    name     string
    functor  PullRopeActionFunc
}

// NewPullRopeFuncAction makes a new PullRopeFuncAction using a text as Named source and the provided function as Caller.
func NewPullRopeFuncAction(
        name      string,
        functor  PullRopeActionFunc) NamedCaller {
    a := new(PullRopeFuncAction)
    a.name = name
    a.functor = functor
    return a
}

// Call executes the pull action.
func (a *PullRopeFuncAction) Call() (res interface{}, err error) {
    err, res = nil, a.functor()
    return
}

// GetName provides access to the pull action name as string return.
func (a *PullRopeFuncAction) GetName() string {
    return fmt.Sprint(a.name)
}

// String implements fmt.Stringer using the pull action name as return.
func (a *PullRopeFuncAction) String() string {
    return fmt.Sprint(a.name)
}








//======================================
// UTILS. logging
//

// StopwatchTrace provides time interval meter access between two instant.
// Usually used as:
//     defer StopwatchTrace("some id text")()
// It returns a function with which emits the elapsed time since was builded.
func StopwatchTrace(slug string) (ender func()) {
    if len(slug) != 0 {
        slug = fmt.Sprintf("[[%v]] ", slug)
    }
    log.Printf("%v🚩Start...", slug)
    start := time.Now()
    return func() {
        log.Printf("%v🏁End. T=%v", slug, time.Since(start))
    }
}






//======================================
// UTILS. randoms
//

type SeededRandConstructor func() *mrand.Rand

// NewTimeSeededSource builds a "math/rand" Source that uses "time" to generate a non-deterministic seed.
//
// Note: this Source still outputs a deterministic sequence based on the seed provided by current "time.Time" UnixNano(), which yields a constantly-changing number.
func NewTimeSeededSource() mrand.Source {
    var seed int64
    // generate constantly-changing seed
    seed = time.Now().UnixNano()
    return mrand.NewSource(seed)
}

// NewTimeSeededRandom is a convenience builder around NewBasicSeededSource().
// It returns a pointer to a "math/rand" Rand that is ready to use.
func NewTimeSeededRandom() *mrand.Rand {
    return mrand.New(NewTimeSeededSource())
}

// NewSecureSeededSource builds a "math/rand" Source that uses "crypto/rand" to generate a non-deterministic seed.
// Note: this Source still outputs a deterministic sequence based on the seed, it's just that the seed is obfuscated.
func NewSecureSeededSource() mrand.Source {
    var seed int64
    // generate obfuscated seed
    binary.Read(crand.Reader, binary.BigEndian, &seed)
    return mrand.NewSource(seed)
}

// NewSecureSeededRandom is a convenience builder around NewSecureSeededSource().
// It returns a pointer to a "math/rand" Rand that is ready to use.
func NewSecureSeededRandom() *mrand.Rand {
    return mrand.New(NewSecureSeededSource())
}








//======================================
// UTILS. sync/atomic counters
//



// A Supplier iface provides access to some generic value.
type Supplier interface {
    Get() interface{}
}

// A Resetter iface is in charge of recover the initial state.
type Resetter interface {
    Reset()
}

// An Incrementer iface allows increment the internal state.
type Incrementer interface {
    Increment()
}

// A Decrementer iface allows decrement the internal state.
type Decrementer interface {
    Decrement()
}

// A SupplierReseter iface embeds the behaviour of both, Supplier and Resetter, managing the resettable state of its generic value.
type SupplierResetter interface {
    Supplier
    Resetter
    GetAndReset() interface{}
    ResetAndGet() interface{}
}

// A SupplierIncrementer iface embeds the behaviour of both, Supplier and Incrementer, managing the incrementable internal state of its generic value.
// In other words, it could be considered as a generic incremental counter.
type SupplierIncrementer interface {
    Supplier
    Incrementer
    GetAndIncrement() interface{}
    IncrementAndGet() interface{}
}

// An IncrementalCounter iface defines the contract to implement by counters with a forwards count capabilities.
type IncrementalCounter SupplierIncrementer

// A SupplierDecrementer iface embeds the behaviour of both, Supplier and Decrementer, managing the decrementable internal state of its generic value.
// In other words, it could be considered as a generic decremental counter.
type SupplierDecrementer interface {
    Supplier
    Decrementer
    GetAndDecrement() interface{}
    DecrementAndGet() interface{}
}
// A DecrementalCounter iface defines the contract to implement by counters with a backwards count capabilities.
type DecrementalCounter SupplierDecrementer

// A Counter iface puts IncrementalCounter and DecrementalCounter capabilities together to manage the internal state of its generic value.
// Then, it's a complete generic counter, either to backwards and forwards.
type Counter interface {
    IncrementalCounter
    DecrementalCounter
}

// A SupplierResetterIncrementer iface embeds the behaviour of SupplierIncrementer but SupplierResetter too, allowing then, to recover the initial state of its generic value.
// In other words, it could be considered as a resettable incremental counter.
type SupplierResetterIncrementer interface {
    SupplierIncrementer
    SupplierResetter
}

// An IncrementalResetCounter iface defines the contract to implement by counters with both capabilities: forwards counting and also resettable itself.
type IncrementalResetCounter SupplierResetterIncrementer

// A SupplierResetterDecrementer iface embeds the behaviour of SupplierDecrementer but SupplierResetter too, allowing then, to recover the initial state of its generic value.
// In other words, can be considered as a resettable decremental counter.
type SupplierResetterDecrementer interface {
    SupplierDecrementer
    SupplierResetter
}

// A DecrementalResetCounter iface defines the contract to implement by counters with both capabilities: backwards counting and also resettable itself.
type DecrementalResetCounter SupplierResetterDecrementer

// A ResetterCounter iface puts Counter and Resetter capabilities together to allow:
//  - Count in both directions, to backwards and to forwards.
//  - Revert it to its initial state.
type ResetterCounter interface {
    Counter
    IncrementalResetCounter
    DecrementalResetCounter
}

/***************** For "INT" type *****************/

// An IntSupplier iface provides access to any value having "int" type.
type IntSupplier interface {
    Supplier
    GetInt() int
}

// An IntSupplierResetter iface embeds the behaviour of both, IntSupplier and Resetter, managing the resettable state of its "int" value.
type IntSupplierResetter interface {
    IntSupplier
    SupplierResetter
    GetIntAndReset() int
    ResetAndGetInt() int
}

// An IntSupplierIncrementer iface embeds the behaviour of both, IntSupplier and Incrementer, managing the incrementable state of its "int" value.
// Due to that, it can be considered as an incremental counter.
type IntSupplierIncrementer interface {
    IntSupplier
    SupplierIncrementer
    IncrementIntN(n int)
    GetIntAndIncrement() int
    GetIntAndIncrementN(n int) int
    IncrementAndGetInt() int
    IncrementAndGetIntN(n int) int
}
// An IntIncrementalCounter iface defines the contract to implement by "int" counters with a forwards count capabilities.
type IntIncrementalCounter IntSupplierIncrementer

// An IntSupplierDecrementer iface embeds the behaviour of both, IntSupplier and Decrementer, managing the decrementable state of its "int" value.
// Due to that, it can be considered as a decremental counter.
type IntSupplierDecrementer interface {
    IntSupplier
    SupplierDecrementer
    DecrementIntN(n int)
    GetIntAndDecrement() int
    GetIntAndDecrementN(n int) int
    DecrementAndGetInt() int
    DecrementAndGetIntN(n int) int
}

// An IntDecrementalCounter iface defines the contract to implement by "int" counters with a backwards count capabilities.
type IntDecrementalCounter IntSupplierDecrementer

// An IntCounter iface puts IntIncrementalCounter and IntDecrementalCounter capabilities together to manage the internal state of its "int" value.
// Then, it's a complete counter, either to backwards and forwards.
type IntCounter interface {
    Counter
    IntIncrementalCounter
    IntDecrementalCounter
}

// A IntSupplierReseterIncrementer iface embeds the behaviour of IntSupplierIncrementer but IntSupplierResetter too, allowing then, to adquire back the initial state of it incrementable "int" value.
// In other words, it could be considered as a resettable incremental counter.
type IntSupplierResetterIncrementer interface {
    IntSupplierIncrementer
    IntSupplierResetter
}

// An IntIncrementalResetCounter iface defines the contract to implement by "int" counters with both capabilities: forwards counting and also resettable itself.
type IntIncrementalResetCounter IntSupplierResetterIncrementer

// An IntSupplierResetterDecrementer iface embeds the behaviour of IntSupplierDecrementer but IntSupplierResetter too, allowing then, to adquire back the initial state of it decrementable "int" value.
// In other words, it could be considered as a resettable decremental counter.
type IntSupplierResetterDecrementer interface {
    IntSupplierDecrementer
    IntSupplierResetter
}

// An IntDecrementalResetCounter iface defines the contract to implement by "int" counters with both capabilities: backwards counting and also resettable itself.
type IntDecrementalResetCounter IntSupplierResetterDecrementer

// A IntResetCounter iface puts together the powefull of IntCounter and Resetter to allow:
//  - Count in both directions, to backwards and to forwards.
//  - Revert it to its "int" initial state.
type IntResetCounter interface {
    IntCounter
    IntIncrementalResetCounter
    IntDecrementalResetCounter
}

type IntCounterConstructor func() IntCounter
type IntCounterConstructorN func(initial int) IntCounter
type IntResetCounterConstructor func() IntResetCounter 
type IntResetCounterConstructorN func(initial int) IntResetCounter



// Helper struct used to serialize/deserialize encapsulated state of reseteable counters.
type intRCounter struct {
    Value    int    `json:"value"    xml:"value"`
    Initial  int    `json:"initial"  xml:"initial"`
}



// IntLockedCounter fully implemments a IntCounter iface by using a sync.Mutex to guard its "int" internal state when the same instance of the counter is shared between several concurrent tasks.
type IntLockedCounter struct {
    lk      sync.Mutex
    value   int        `json:"value"  xml:"value"`
}

// NewIntLockedCounter makes a new IntLockedCounter instance initialized to zero.
func NewIntLockedCounter() IntCounter {
    return &IntLockedCounter{ value: 0 }
}

// NewIntLockedCounterN makes a new IntLockedCounter instance initialized to some value.
func NewIntLockedCounterN(initial int) IntCounter {
    c := new(IntLockedCounter)
    c.value   = initial
    return c
}

// String implements fmt.Stringer to get a string representation of the counter's current value.
func (c *IntLockedCounter) String() string {
    value := c.get()
    return strconv.Itoa(value)
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the counter internal value as JSON number.
func (c *IntLockedCounter) MarshalJSON() ([]byte, error) {
    return json.Marshal(c.get())
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the counter "int" value.
func (c *IntLockedCounter) UnmarshalJSON(b []byte) (err error) {
    var v int
    if err = json.Unmarshal(b, &v); err == nil {
        c.set(v)
    }
    return
}

func (c *IntLockedCounter) Get() interface{} {
    return c.get()
}

func (c *IntLockedCounter) GetInt() int {
    return c.get()
}

func (c *IntLockedCounter) Increment() {
    c.add(1)
}

func (c *IntLockedCounter) IncrementIntN(n int) {
    checkSyncCounterIncInt(n)
    c.add(n)
}

func (c *IntLockedCounter) GetAndIncrement() interface{} {
    return c.getadd(1)
}

func (c *IntLockedCounter) GetIntAndIncrement() int {
    return c.getadd(1)
}

func (c *IntLockedCounter) GetIntAndIncrementN(n int) int {
    checkSyncCounterIncInt(n)
    return c.getadd(n)
}

func (c *IntLockedCounter) IncrementAndGet() interface{} {
    return c.addget(1)
}

func (c *IntLockedCounter) IncrementAndGetInt() int {
    return c.addget(1)
}

func (c *IntLockedCounter) IncrementAndGetIntN(n int) int {
    checkSyncCounterIncInt(n)
    return c.addget(n)
}

func (c *IntLockedCounter) Decrement() {
    c.add(-1)
}

func (c *IntLockedCounter) DecrementIntN(n int) {
    checkSyncCounterDecInt(n)
    c.add(-n)
}

func (c *IntLockedCounter) GetAndDecrement() interface{} {
    return c.getadd(-1)
}

func (c *IntLockedCounter) GetIntAndDecrement() int {
    return c.getadd(-1)
}

func (c *IntLockedCounter) GetIntAndDecrementN(n int) int {
    checkSyncCounterDecInt(n)
    return c.getadd(-n)
}

func (c *IntLockedCounter) DecrementAndGet() interface{} {
    return c.addget(-1)
}

func (c *IntLockedCounter) DecrementAndGetInt() int {
    return c.addget(-1)
}

func (c *IntLockedCounter) DecrementAndGetIntN(n int) int {
    checkSyncCounterDecInt(n)
    return c.addget(-n)
}

// DoWithValue implements functional.ValueConsumer allowing access to the current value inside its concurrent safe context.
func (c *IntLockedCounter) DoWithValue(consumer ConsumerFunc) {
    c.lk.Lock()
    defer c.lk.Unlock()
    // execute callback in this safe context
    consumer(c.value)
}

// DoWithIntValue implements functional.IntValueConsumer allowing access to the current value inside its concurrent safe context.
func (c *IntLockedCounter) DoWithIntValue(consumer IntConsumerFunc) {
    c.lk.Lock()
    defer c.lk.Unlock()
    // execute callback in this safe context
    consumer(c.value)
}

func (c *IntLockedCounter) set(n int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    c.value = n
}

func (c *IntLockedCounter) get() int {
    c.lk.Lock()
    defer c.lk.Unlock()
    return c.value
}

func (c *IntLockedCounter) add(n int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    c.value += n
}

func (c *IntLockedCounter) addget(n int) (value int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    c.value += n
    value = c.value
    return
}

func (c *IntLockedCounter) getadd(n int) (value int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    value = c.value
    c.value += n
    return
}



// IntLockedResetCounter is a counter of the IntResetCounter family. It extends an IntLockedCounter with Resetter capabilities, this is, it can revert to its "int" initial state with which it was created.
type IntLockedResetCounter struct {
    IntLockedCounter
    initial    int
}

// NewIntLockedResetCounter makes a new IntLockedResetCounter instance initialized to zero.
func NewIntLockedResetCounter() IntResetCounter {
    c := new(IntLockedResetCounter)
    c.value   = 0
    c.initial = 0
    return c
}

// NewIntLockedResetCounter makes a new IntLockedResetCounter instance initialized to some value.
func NewIntLockedResetCounterN(initial int) IntResetCounter {
    c := new(IntLockedResetCounter)
    c.value   = initial
    c.initial = initial
    return c
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the counter internal state (value and initial "int"s) to JSON.
func (c *IntLockedResetCounter) MarshalJSON() ([]byte, error) {
    // Delegating into other struct is needed because values
    // are not public visible to json package.
    var v = intRCounter{
        Initial: c.initial,
        Value: c.get(),
        }
    return json.Marshal(v)
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the counter "int" value and initial state.
func (c *IntLockedResetCounter) UnmarshalJSON(b []byte) (err error) {
    // Delegating into other struct is needed because values
    // are not public visible to json package.
    var v intRCounter
    if err = json.Unmarshal(b, &v); err == nil {
        c.initial = v.Initial
        c.set(v.Value)
    }
    return
}

// Reset sets the initial value as current value.
func (c *IntLockedResetCounter) Reset() {
    c.set(c.initial)
}

// GetAndReset sets the initial value as current value returning this last state.
func (c *IntLockedResetCounter) GetAndReset() interface{} {
    return c.getreset()
}

// GetIntAndReset sets the initial value as current value returning this last state.
func (c *IntLockedResetCounter) GetIntAndReset() int {
    return c.getreset()
}

// ResetAndGet sets the initial value as current value returning the initial value.
func (c *IntLockedResetCounter) ResetAndGet() interface{} {
    return c.resetget()
}

// ResetAndGetInt sets the initial value as current value returning the initial value.
func (c *IntLockedResetCounter) ResetAndGetInt() int {
    return c.resetget()
}

func (c *IntLockedResetCounter) getreset() (value int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    value = c.value
    c.value = c.initial
    return
}

func (c *IntLockedResetCounter) resetget() (value int) {
    c.lk.Lock()
    defer c.lk.Unlock()
    value = c.initial
    c.value = value
    return
}

func checkSyncCounterIncInt(n int) {
    if n < 1 {
        panic("sync/counter: invalid increment")
    }
}
func checkSyncCounterDecInt(n int) {
    if n < 1 {
        panic("sync/counter: invalid decrement")
    }
}






// IntAtomicCounter fully implemments a IntCounter iface by using atomic operations to guard its "int" internal state when the same instance of the counter is shared between several concurrent tasks.
type IntAtomicCounter struct {
    value    unsafe.Pointer  // *int
}

// NewIntAtomicCounter makes a new IntAtomicCounter instance initialized to zero.
func NewIntAtomicCounter() IntCounter {
    return NewAtomicCounterIntN(0)
}

// NewAtomicCounterIntN makes a new IntAtomicCounter instance initialized to some value.
func NewAtomicCounterIntN(initial int) IntCounter {
    c := new(IntAtomicCounter)
    c.value = unsafe.Pointer(&initial)
    return c
}

// String implements fmt.Stringer to get a string representation of the counter's current value.
func (c *IntAtomicCounter) String() string {
    value := c.get()
    return strconv.Itoa(value)
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the counter internal value as JSON number.
func (c *IntAtomicCounter) MarshalJSON() ([]byte, error) {
    return json.Marshal(c.get())
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the counter "int" value.
func (c *IntAtomicCounter) UnmarshalJSON(b []byte) (err error) {
    var v int
    if err = json.Unmarshal(b, &v); err == nil {
        c.set(v)
    }
    return
}

func (c *IntAtomicCounter) Get() interface{} {
    return c.get()
}

func (c *IntAtomicCounter) GetInt() int {
    return c.get()
}

func (c *IntAtomicCounter) Increment() {
    c.add(1)
}

func (c *IntAtomicCounter) IncrementIntN(n int) {
    checkAtomicCounterIncInt(n)
    c.add(n)
}

func (c *IntAtomicCounter) GetAndIncrement() interface{} {
    return c.getadd(1)
}

func (c *IntAtomicCounter) GetIntAndIncrement() int {
    return c.getadd(1)
}

func (c *IntAtomicCounter) GetIntAndIncrementN(n int) int {
    checkAtomicCounterIncInt(n)
    return c.getadd(n)
}

func (c *IntAtomicCounter) IncrementAndGet() interface{} {
    return c.addget(1)
}

func (c *IntAtomicCounter) IncrementAndGetInt() int {
    return c.addget(1)
}

func (c *IntAtomicCounter) IncrementAndGetIntN(n int) int {
    checkAtomicCounterIncInt(n)
    return c.addget(n)
}

func (c *IntAtomicCounter) Decrement() {
    c.add(-1)
}

func (c *IntAtomicCounter) DecrementIntN(n int) {
    checkAtomicCounterDecInt(n)
    c.add(-n)
}

func (c *IntAtomicCounter) GetAndDecrement() interface{} {
    return c.getadd(-1)
}

func (c *IntAtomicCounter) GetIntAndDecrement() int {
    return c.getadd(-1)
}

func (c *IntAtomicCounter) GetIntAndDecrementN(n int) int {
    checkAtomicCounterDecInt(n)
    return c.getadd(-n)
}

func (c *IntAtomicCounter) DecrementAndGet() interface{} {
    return c.addget(-1)
}

func (c *IntAtomicCounter) DecrementAndGetInt() int {
    return c.addget(-1)
}

func (c *IntAtomicCounter) DecrementAndGetIntN(n int) int {
    checkSyncCounterDecInt(n)
    return c.addget(-n)
}

func (c *IntAtomicCounter) DoWithValue(consumer ConsumerFunc) {
    value := c.get()
    // execute callback in this safe context
    consumer(value)
}

func (c *IntAtomicCounter) DoWithIntValue(consumer IntConsumerFunc) {
    value := c.get()
    // execute callback in this safe context
    consumer(value)
}

func (c *IntAtomicCounter) get() (value int) {
    // load current value
    p := atomic.LoadPointer(&c.value)
    // get Pointer internal value
    value = *(*int)(p)
    return
}

func (c *IntAtomicCounter) set(n int) {
    p := unsafe.Pointer(&n)
    atomic.StorePointer(&c.value, p)
}

func (c *IntAtomicCounter) add(delta int) {
    for { // CAS (Compare-And-Swap) loop
        // load current value
        old := atomic.LoadPointer(&c.value)
        // get Pointer internal value
        value := *(*int)(old)
        // apply delta
        new := value + delta
        // check until really setted
        if atomic.CompareAndSwapPointer(&c.value, old, unsafe.Pointer(&new)) {
            return
        }
    }
}

func (c *IntAtomicCounter) addget(delta int) int {
    for { // CAS (Compare-And-Swap) loop
        // load current value
        old := atomic.LoadPointer(&c.value)
        // get Pointer internal value
        value := *(*int)(old)
        // apply delta
        new := value + delta
        // check until really setted
        if atomic.CompareAndSwapPointer(&c.value, old, unsafe.Pointer(&new)) {
            return new
        }
    }
}

func (c *IntAtomicCounter) getadd(delta int) int {
    for { // CAS (Compare-And-Swap) loop
        // load current value
        old := atomic.LoadPointer(&c.value)
        // get Pointer internal value
        value := *(*int)(old)
        // apply delta
        new := value + delta
        // check until really setted
        if atomic.CompareAndSwapPointer(&c.value, old, unsafe.Pointer(&new)) {
            return value
        }
    }
}


// IntAtomicResetCounter is a counter of the IntResetCounter family. It extends an IntAtomicCounter with Resetter capabilities, this is, it can revert to its "int" initial state with which it was created.
type IntAtomicResetCounter struct {
    IntAtomicCounter
    initial    int
}

// NewIntAtomicResetCounter makes a new IntAtomicResetCounter instance initialized to zero.
func NewIntAtomicResetCounter() IntResetCounter {
    return NewIntAtomicResetCounterN(0)
}

// NewIntAtomicResetCounterN makes a new IntAtomicResetCounter instance initialized to some value.
func NewIntAtomicResetCounterN(initial int) IntResetCounter {
    c := new(IntAtomicResetCounter)
    c.value   = unsafe.Pointer(&initial)
    c.initial = initial
    return c
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the counter internal state (value and initial "int"s) to JSON.
func (c *IntAtomicResetCounter) MarshalJSON() ([]byte, error) {
    // Delegating into other struct is needed because values
    // are unsafe.Pointers and/or also not public visible to json package.
    var v = intRCounter{
        Initial: c.initial,
        Value: c.get(),
        }
    return json.Marshal(v)
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the counter "int" value and initial state.
func (c *IntAtomicResetCounter) UnmarshalJSON(b []byte) (err error) {
    // Delegating into other struct is needed because values
    // are unsafe.Pointers and/or also not public visible to json package.
    var v intRCounter
    if err = json.Unmarshal(b, &v); err == nil {
        c.initial = v.Initial
        c.set(v.Value)
    }
    return
}

// Reset sets the initial value as current value.
func (c *IntAtomicResetCounter) Reset() {
    c.set(c.initial)
}

// GetAndReset sets the initial value as current value returning this last state.
func (c *IntAtomicResetCounter) GetAndReset() interface{} {
    return c.getreset()
}

// GetIntAndReset sets the initial value as current value returning this last state.
func (c *IntAtomicResetCounter) GetIntAndReset() int {
    return c.getreset()
}

// ResetAndGet sets the initial value as current value returning the initial value.
func (c *IntAtomicResetCounter) ResetAndGet() interface{} {
    return c.resetget()
}

// ResetAndGetInt sets the initial value as current value returning the initial value.
func (c *IntAtomicResetCounter) ResetAndGetInt() int {
    return c.resetget()
}

func (c *IntAtomicResetCounter) getreset() (value int) {
    initialp := unsafe.Pointer(&c.initial)
    for { // CAS (Compare-And-Swap) loop
        // load current value
        valuep := atomic.LoadPointer(&c.value)
        // get Pointer internal value
        value = *(*int)(valuep)
        // check until really setted
        if atomic.CompareAndSwapPointer(&c.value, valuep, initialp) {
            return
        }
    }
}

func (c *IntAtomicResetCounter) resetget() (value int) {
    value = c.initial
    c.set(c.initial)
    return
}

func checkAtomicCounterIncInt(n int) {
    if n < 1 {
        panic("sync/atomic: invalid increment")
    }
}
func checkAtomicCounterDecInt(n int) {
    if n < 1 {
        panic("sync/atomic: invalid decrement")
    }
}



// IntCounterRandomFactory is a method factory that generates an aleatory IntCounter instance from each one registered in the provider of same type.
func IntCounterRandomFactory(rnd *mrand.Rand) func() IntCounter {
    n := rnd.Int() % int(maxIntCounters)
    return IntCounterProvider(n).New
}

// IntCounterProvider holds the register of all IntCounter data types defined in the package and allowing create instances of each one identified by an int number.
type IntCounterProvider int
const (
                   // iota generate ints from zero.
                   // Is reseted in each (...) group
    Locked         IntCounterProvider = iota
    Atomic
    LockedReset
    AtomicReset
    maxIntCounters
)

// String implements fmt.Stringer to get a string representation of the counter struct that represent the enumerated constants defined by this provider.
func (p IntCounterProvider) String() string {
    if p < 0 || p >= maxIntCounters {
        return "Unknown IntCounterProvider constant #" + strconv.Itoa(int(p))
    }
    return [...]string{
        "IntLockedCounter",
        "IntAtomicCounter",
        "IntLockedResetCounter",
        "IntAtomicResetCounter",
    }[p]
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the enumerated constant as JSON number.
func (p IntCounterProvider) MarshalJSON() (b []byte, err error) {
    if v := int(p); p < 0 || p >= maxIntCounters {
        b, err = nil, fmt.Errorf("json: unknown IntCounterProvider constant: %d", v)
    } else {
        b, err = json.Marshal(v)
    }
    return
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON into the value matching the enumerated constants defined by this provider.
func (p *IntCounterProvider) UnmarshalJSON(b []byte) (err error) {
    var v int
    if err = json.Unmarshal(b, &v); err != nil {
        return
    }
    if v < 0 || v >= int(maxIntCounters) {
        err = fmt.Errorf("json: unknown IntCounterProvider constant: %d", v)
    } else {
        *p, err = IntCounterProvider(v), nil
    }
    return
}

// New makes a new instance of the IntCounter that it's related with the provider registry constant.
func (p IntCounterProvider) New() (inst IntCounter) {
    switch p {
        case Locked:
            inst = NewIntLockedCounter()
        case Atomic:
            inst = NewIntAtomicCounter()
        case LockedReset:
            inst = NewIntLockedResetCounter()
        case AtomicReset:
            inst = NewIntAtomicResetCounter()
        default:
            panic("counters: requested IntCounter #" + strconv.Itoa(int(p)) + " is unavailable")
    }
    return
}






//======================================
// UTILS. functional
//

type ConsumerFunc        func(value interface{})
type PointerConsumerFunc func(value unsafe.Pointer)
type ByteConsumerFunc    func(value byte)
type BoolConsumerFunc    func(value bool)
type IntConsumerFunc     func(value int)
type Int8ConsumerFunc    func(value int8)
type Int16ConsumerFunc   func(value int16)
type Int32ConsumerFunc   func(value int32)
type Int64ConsumerFunc   func(value int64)
type UintConsumerFunc    func(value uint)
type Uint8ConsumerFunc   func(value uint8)
type Uint16ConsumerFunc  func(value uint16)
type Uint32ConsumerFunc  func(value uint32)
type Uint64ConsumerFunc  func(value uint64)
type UintptrConsumerFunc func(value uintptr)
type StringConsumerFunc  func(value string)
type RuneConsumerFunc    func(value rune)


// A ValueConsumer iface allows access generic values provided by a callback consumer function that receive that value as argument.
type ValueConsumer interface {
    DoWithValue(f ConsumerFunc)
}
// A IntValueConsumer iface allows access "int" values provided by a callback consumer function that receive that value as argument.
type IntValueConsumer interface {
    DoWithIntValue(f IntConsumerFunc)
}
