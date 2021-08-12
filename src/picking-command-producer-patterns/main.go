/*
🐹 Crazy and unpredictable calculator

TOPIC: OOP, Patterns and basic concurrency in Go!

Design Patterns you can see coded here:
- Closure
- Composition
- Asynchronous Producer/Consumer
- Command
- Strategy
- Adapter
- Decorator
- Singleton
- Memoizer??
- Template/Supports partial structs
- Multiple inheritance/embedded

Clean code / SOLID:
- Interface segregation
- OCP. Open-Close Principle
- NewXXX method signatures return ifaces
- Composition vs Inheritance
- TDD ready


@author davorpatech
@since 2021-05-07
*/
package main

import (
    "errors"
    "fmt"
    "io"
    "log"
    "math"
    "math/rand"
    "os"
    "reflect"
    "strconv"
    "strings"
    "sync"
    "time"
)

//======================================
// PROGRAM
//

// Global config

const delay time.Duration = 25 * time.Millisecond
const queueSize int = 30
//⚠ SL max exec time: 1s > queueSize*delay

var (
    logerr    *log.Logger
)

func init() { // init package
    const flags = log.Ltime | log.Lmicroseconds | log.LUTC | log.Lmsgprefix
    // Configure default logger
    log.SetFlags(flags)
    log.SetOutput(os.Stdout)
    log.SetPrefix(":app: ")
    // configure error logger
    logerr = log.New(os.Stdout, ":err: ", flags)
}


func main() {
    defer TraceStartEnd(time.Now())()
    printSplash()
    
    var inum float64 = input()
    
    // Operation registry
    //@TODO: avoid mutate "n" on closures??
    //@TODO provide other command type
    n, mem1 := inum, inum
    ops := []Command{
        NewCmdFunc("Add.3",
            func(){ n += 3.0 }),
        NewCmdFunc("Substract.1",
            func(){ n -= 1.0 }),
        NewCmdFunc("Doubling",
            func(){ n *= 2.0 }),
        NewCmdFunc("Halfing",
            func(){ n /= 2.0 }),
        NewCmdFunc("Multiply.Self",
            func(){ n *= n }),
        NewCmdFunc("Root.2",
            func(){ n = math.Sqrt(n) }),
        NewCmdFunc("Root.3",
            func(){ n = math.Cbrt(n) }),
        NewCmdFunc("Pow.3",
            func(){ n = math.Pow(n, 3) }),
        NewCmdFunc("Abs",
            func(){ n = math.Abs(n) }),
        NewCmdFunc("Ceiling",
            func(){ n = math.Ceil(n) }),
        NewCmdFunc("Flooring",
            func(){ n = math.Floor(n) }),
        NewCmdFunc("Round",
            func(){ n = math.Round(n) }),
        NewCmdFunc("Round.2Even",
            func(){ n = math.RoundToEven(n) }),
        NewCmdFunc("Exp.E",
            func(){ n = math.Exp(n) }),
        NewCmdFunc("Exp.2",
            func(){ n = math.Exp2(n) }),
        NewCmdFunc("Reset.NaN/Inf",
            func(){ if isNanInf(n) { n = inum } }),
        NewCmdFunc("Undo.NaN/Inf",
            func(){ if isNanInf(n) { n = mem1 } }),
        NewCmdFunc("Log.N",
            func(){ n = math.Log(n) }),
        NewCmdFunc("Log.10",
            func(){ n = math.Log10(n) }),
        NewCmdFunc("Remaind.69",
            func(){ n = math.Mod(n, 69) }),
        NewCmdFunc("Extract.Fract",
            func(){ _, n = math.Modf(n) }),
    }
    
    // create Producer by composition
    var prod CommandProducer
    prod = NewPickingCommandProducer(
        NewSliceCommandPicker(ops), 
        true)
    // create unbuffered receiver channel
    recv := make(chan Command)
    
    // fill channel in a routine
    go prod.ProduceN(recv, queueSize)

    // wait for data in channel & drain it
    for op := range recv {
        mem1 = n
        // polymorphism Command=>Executable
        // mutate "n"
        op.Execute()
        // Output
        fmt.Printf("%13[1]s(%20.10[2]f) =\n%[4]s%20.10[3]f\n\n", 
            op, //polymorphism fmt.Stringer
            mem1, n,
            strings.Repeat(" ", 14))
        // Pause
        time.Sleep(delay)
    }
    
    // wait for data
    <- recv
    
    //Test/Check iface implementation/polimorphism
    ops[0].(Wrappeer).Unwrap()
    ops[0].(Delegater).GetDelegateObj()
    ops[0].(Delegater).GetDelegateType()
    ops[0].(fmt.Stringer).String()
    ops[0].(ComponentNameer).ComponentName()
    prod.(Delegater).GetDelegateObj()
    prod.(Delegater).GetDelegateType()
}


func printSplash() {
    w := 39
    fmt.Printf(` %[1]v
 %[2]v
 |%[3]v|
<|%[4]v|>
 |%[3]v|
 %[2]v
`,
        CenterStr("_.o•°~@/¤\\@~°•o._", w),
        strings.Repeat("=", w),
        strings.Repeat(" ", w-2),
        CenterStr("🐹 The unpredicted calculator 🐹", w-2*2),    //2=>emojis
        )
}


func input() (num float64) {
    num = 8.0    // default
    fmt.Print("Input number: ")
    defer func() {
        fmt.Println()
        fmt.Println()
    }()
    if n, err := fmt.Scan(&num); n != 1 {
        switch {
            case err == io.EOF, err == io.ErrUnexpectedEOF, err.Error() == "unexpected newline":
                // use default
                //break
            case err == strconv.ErrRange, errors.Is(err, strconv.ErrRange):
                Dien(1, "Input number out of range.")
            case err == strconv.ErrSyntax, errors.Is(err, strconv.ErrSyntax):
                Dien(1, "Input not is a number.")
            case errors.Is(err, strconv.ErrRange):
                log.Println("Input number overflow.")
            default:
                Dien(2, "Input number: ", err)
        }
        
        // use default
        log.Println("ℹ", "Using default value as input: ", num)
    }
    return
}





//======================================
// MODEL
//


// Named Component iface.
type ComponentNameer interface {
    // embed fmt.Stringer iface
    fmt.Stringer
    
    ComponentName() string
}

// Wrappable iface.
type Wrappeer interface {
    Unwrap() interface{}
}

// Delegater iface as Wrappeer adapter.
type Delegater interface {
    // embed super iface (adaptes)
    Wrappeer
    
    GetDelegateObj() interface{}
    GetDelegateType() string
}

// Executable iface.
type Executable interface {
    Execute()
}

// Command iface.
type Command interface {
    // embed super ifaces (composition)
    ComponentNameer
    Executable
}

// CommandProducer iface.
// It is in charge of produce commands asynchronously, writting them in an incomming channel.
type CommandProducer interface {
    ProduceN(recv chan <- Command, n int)
}

// CommandPicker iface.
// It acts as Command selector.
type CommandPicker interface {
    PickCmd() Command
}



// Abstract type that gives common support to a Delegater iface.
// Subtypes have to complete iface contract implementing the Delegater#GetDelegateObj() method and set a reference to itself.
// See also InitDelegaterSupport.
type DelegateSupport struct {
    // partially embed Delegater iface
    Delegater
}
// Unwrap method implements Wrappeer iface using Delegater#GetDelegateObj.
func (self DelegateSupport) Unwrap() interface{} {
    return self.GetDelegateObj()
}
// GetDelegateType method implements Delegater iface obtaining the typeof(GetDelegateObj).
func (self DelegateSupport) GetDelegateType() (t string) {
    obj := self.GetDelegateObj()
    // both does the same
    t = reflect.TypeOf(obj).String()
    t = fmt.Sprintf("%T", obj)
    return
}

var tDelegaterSupp = reflect.TypeOf(DelegateSupport{})

// InitDelegaterSupport inntrospect searching right Delegater field and complete Delegater abstraction allowing lookup "obj" in DelegaterSupport#GetDelegateType() method using the target implementation as source
func InitDelegaterSupport(target Delegater) Delegater {
    if target == nil {
        panic("core: Delegater to initialize its support is not provided.")
    }
    
    // introspect unwrapping *ptr if needed
    Vof := reflect.ValueOf
    v := reflect.Indirect(Vof(target))
    if tv := v.Type(); tv == tDelegaterSupp || tv.AssignableTo(tDelegaterSupp) {
        panic("core: DelegateSupport cannot be Delegater of itself.")
    }
    
    // Set emmbed field using reflection
    if f := v.FieldByName("DelegateSupport"); f.IsValid() {
        var supporter Delegater
        if f.Kind() == reflect.Ptr {
            // as pointer
            supporter = &DelegateSupport{
                Delegater: target }
        } else {
            // as struct
            supporter = DelegateSupport{
                Delegater: target }
        }
        f.Set(Vof(supporter))
    }
    
    /* static way
    // this complete Delegater abstraction allowing lookup "obj" in DelegateSupport#GetDelegateType() method using the target implementation as source
    target.DelegateSupport = DelegateSupport{
        
        Delegater: target }
    */
    return target
}


// Abstract type that acts as common  Command, providing base info to ComponentNameer and then fmt.Stringer.
// Subtypes have to implement Executable to complete the Command iface contract.
type BaseCommand struct {
    // partially embed Command iface
    Command
    // self attributes
    name string    `json:"name"  xml:"name"`
}
// ComponentName method implements ComponentNameer iface and then, also part of Command.
func (self *BaseCommand) ComponentName() string {
    return self.name
}
// String method implements fmt.Stringer iface using #ComponentName().
func (self *BaseCommand) String() string {
    return self.ComponentName()
}


// NewCmdFunc method builds a new Command using a lambda as Executable source.
func NewCmdFunc(name string, fn func()) Command {
    if IsBlankStr(name) {
        panic("command: name not provided to NewCmdFunc.")
    }
    if fn == nil {
        panic("command: lambda function not provided to NewCmdFunc.")
    }
    this := &funcCmdAdapter{
            fn: fn,
        }
    //=== provide embbed info
    // set BaseCommand name when no pointer
//    this.name = name
    // set info when embed pointer
    this.BaseCommand = &BaseCommand{
        Command: this,
        name: name }
    InitDelegaterSupport(this)
    // return instance
    return this
}

// LambdaCommandAdapter acts as concrete class implementing Command & Executable & Delegater. It use a lambda as delegated Executable source.
type funcCmdAdapter struct {
    // embbed super structs (multi-inheritance)
    *BaseCommand    // as ptr (r/w access)
    *DelegateSupport
    // self attributes
    fn func()        // lambda source
}
// Execute method implements Executable iface and then Command, closing BaseCommand abstraction.
func (self *funcCmdAdapter) Execute() {
    self.fn()
}
// GetDelegateObj method implements Delegater iface, closing DelegateSupport abstraction.
// It acts as accessor to the lambda function itself.
func (self *funcCmdAdapter) GetDelegateObj() interface{} {
    return self.fn
}



// NewSliceCommandPicker method builds a new CommandPicker with provided command slice as picking source
func NewSliceCommandPicker(cmds []Command) CommandPicker {
    if len(cmds) < 1 {
        panic("command: Empty CommandPicker slice source.")
    }
    this := &sliceCommandPicker{
            cmds: cmds,
        }
    this.init()
    return this
}

// SliceCommandPicker acts as concrete class implementing CommandPicker having a Command slice as picking source.
type sliceCommandPicker struct {
    // singleton initializer context
    once sync.Once
    m    sync.Mutex    // thread-safe monitor
    rnd  *rand.Rand    // random instance
    // self attributes
    cmds []Command     // slice data
}
func (self *sliceCommandPicker) init() {
    // lock concurrent access
    self.m.Lock()
    defer self.m.Unlock()
    // Init random seeder only once.
    self.once.Do(func() {
        // Seeding with the same value results in the same random sequence each run.
        // For different numbers, seed with a different value, such as
        // time.Now().UnixNano(), which yields a constantly-changing number.
        // for a secure seed use "crypto/rand" instead "time"
        seeder := rand.NewSource(time.Now().UnixNano())
        self.rnd = rand.New(seeder)
    })
}
// PickCmd method implements CommandPicker iface.
// It retrieves a Command from the slice using a random algorithm.
func (self *sliceCommandPicker) PickCmd() Command {
    // lock concurrent access
    self.m.Lock()
    defer self.m.Unlock()
    // compute
    siz := len(self.cmds)
    rint:= self.rnd.Intn(siz)
    cmd := self.cmds[rint]
    return cmd
}


// NewPickingCommandProducer method builds a new CommandProducer with provided picker as picking strategy
func NewPickingCommandProducer(cpkr CommandPicker, closeReceiver bool) CommandProducer {
    if cpkr == nil {
        panic("producer: CommandPicker not provided.")
    }
    var this CommandProducer
    this = &pickingCommandProducer{
            cpkr: cpkr,
        }
    if closeReceiver { // decorate
        this = NewClosingCommandProducerDecorator(this)
    }
    // return instance
    return this
}

// PickingCommandProducer acts as concrete composite class implementing CommandProducer iface taking produced commands from any CommandPicker.
type pickingCommandProducer struct {
    // self attributes
    cpkr CommandPicker  // picker strategy
}
// ProduceN method implements CommandProducer interface.
// It takes as its parameters a Command channel and the number describing how many picked commands it should write into this channel.
func (self *pickingCommandProducer) ProduceN(recv chan <- Command, n int) {
    // PROMOTED AS Decorator. It flexibilize composition & TDD
    // ensure close channel when gofunc ends, we know fixed length!
    // Combined with "for x := range chan" iteration, marks the end of that loop
    //defer close(recv)
    
    // generate n-commands
    for i := 0; i < n; i++ {
        cmd := self.cpkr.PickCmd()
        // write / notify
        recv <- cmd
    }
}


// NewClosingCommandProducerDecorator method builds a new CommandProducer decorator that closes underliying channel after produce Commands.
func NewClosingCommandProducerDecorator(delegate CommandProducer) CommandProducer {
    if delegate == nil {
        panic("producer: CommandProducer delegate to decorate not provided.")
    }
    // type assertion to avoid double wrap
    // using a type switch
    switch delegate.(type) {
        case *closingCommandProducer:
             panic("producer: CommandProducer delegate to decorate cannot be decorated itself as *pointer.")
        case closingCommandProducer:
             panic("producer: CommandProducer delegate to decorate cannot be decorated itself.")
    }
    this := &closingCommandProducer{
            delegate: delegate,
        }
    InitDelegaterSupport(this)
    // return instance
    return this
}

// ClosingCommandProducer acts as concrete CommandProducer iface decorator that closes underliying channel after produce Commands.
type closingCommandProducer struct {
    // embbed super structs
    DelegateSupport
    // self attributes
    delegate CommandProducer    // to decorate
}
// ProduceN method implements CommandProducer interface.
// It takes as its parameters a Command channel and the number describing how many picked commands it should write into this channel.
// This channel is closed when we're done producing values.
func (self closingCommandProducer) ProduceN(recv chan <- Command, n int) {
    // ensure close channel when gofunc ends, we know fixed length!
    // Combined with "for x := range chan" iteration, marks the end of that loop
    defer close(recv)
    // call decorated delegate
    self.delegate.ProduceN(recv, n)
}
// GetDelegateObj method implements Delegater iface, closing DelegateSupport abstraction. It acts as accessor to the underliying decorated CommandProducer itself.
func (self closingCommandProducer) GetDelegateObj() interface{} {
    return self.delegate
}





//======================================
// UTILS. Error interoperability
//

// EDie halts process with error as the output message.
// Like log.Fatalln(...)
func EDie(e error) {
    EDien(1, e)
}
// EDien halts process with specific code and error as the output message.
// Like log.Fatalln(...)
func EDien(code int, e error) {
    Dien(code, "⚠", e)
}
// Dien halts process with specific code and a variadic arguments as an error output message.
// Like log.Fatalln(...)
func Dien(code int, args ...interface{}) {
    Eprintln(args...)
    os.Exit(code)
}
// Dienf halts process with specific code and a variadic arguments as an error output message formatted according to the desired pattern.
// Like log.Fatalf(...)
func Dienf(code int, format string, args ...interface{}) {
    s := fmt.Sprintf(format, args...)
    Eprintln("⚠", s)
    os.Exit(code)
}

// Eprintln prints the variadic arguments into os.Stderr writer.
func Eprintln(args ...interface{}) (n int, err error) {
    n, err = fmt.Fprintln(os.Stderr, args...)
    return
}



//======================================
// UTILS. Math
//


// isNanInf returns true if a number is NaN or -Inf/+Inf.
func isNanInf(f float64) bool {
    return math.IsNaN(f) || isInf(f)
}

// IsInf returns true if a number is -Inf/+Inf
func isInf(f float64) bool {
    // 0 => +Inf or -Inf
    return math.IsInf(f, 0)
}



//======================================
// UTILS. time
//

// TraceStartEnd emits a log trace at begin and end with the elapsed time on invoke returned function.
// Commonly used as... defer TraceStartEnd(time.Now())()
func TraceStartEnd(start time.Time) (ender func()) {
    log.Println("ℹ🚩 Start...")
    return func() {
        log.Println("ℹ🏁 End.", time.Since(start))
    }
}



//======================================
// UTILS. String
//

// StrCenter tries to center a string to a given width using whitespace as filling characters.
// It preserves codepoints joined when distributable width is computed.
func CenterStr(s string, w int) (r string) {
    if len(s) == 0 {
        r = strings.Repeat(" ", w)
        return
    }
    scp := []rune(s)
//    r = fmt.Sprintf(Sprintf("%%-%ds", w/2), fmt.Sprintf(fmt.Sprintf("%%%ds", w/2),s))
//    r = fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w + len(s))/2, s))
    r = fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w + len(scp))/2, s))
    return
}

// IsBlankStr returns true if string to check has only whitespaces.
func IsBlankStr(s string) bool {
    return IsEmptyStr(strings.TrimSpace(s))
}
// IsEmptyStr returns true if string to check hasn't length.
func IsEmptyStr(s string) bool {
    return len(s) == 0
}
