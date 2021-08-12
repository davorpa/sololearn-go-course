/*
Practising type alias, implementing interfaces, and providing wrappers to get the same behaviour using inheritance that by composition (wrapper)

@author davorpatech
@since 2021-06-10
*/
package main

import "encoding/json"
import "fmt"
import "log"
import "os"
import "strconv"
import "sync"


// Inter iface contract
type Inter interface {
    GetInt() int
    SetInt(value int)
}

func IsInter(value interface{}) bool {
    // safe type casting
    if _, ok := value.(Inter); ok {
        return true
    }
    return false
}

type Incrementer interface {
    Inc() Incrementer
}

func IsIncrementer(value interface{}) bool {
    // safe type casting
    if _, ok := value.(Incrementer); ok {
        return true
    }
    return false
}

type Decrementer interface {
    Dec() Decrementer
}

func IsDecrementer(value interface{}) bool {
    // safe type casting
    if _, ok := value.(Decrementer); ok {
        return true
    }
    return false
}


type InterOps interface {
    Incrementer
    Decrementer
    // ....
}



// An BoxedInt is an Inter arround a primitive "int" type.
type BoxedInt struct {
    value    int
}

/*    Constructors    */

func Int(value int) Inter {
    n := new(BoxedInt)
    n.value = value
    return n
}

/*    implement Inter    */

func (n *BoxedInt) SetInt(value int) {
    log.Println("BoxedInt.Set")
    n.value = value
}

func (n *BoxedInt) GetInt() int {
    log.Println("BoxedInt.Get")
    return n.value
}

/*    implement Incrementer    */

func (n *BoxedInt) Inc() Incrementer {
    n.SetInt(n.GetInt() + 1)
    return n
}

/*    implement Decrementer    */

func (n *BoxedInt) Dec() Decrementer {
    log.Println("BoxedInt.Dec")
    n.SetInt(n.GetInt() - 1)
    return n
}

/*    implement fmt.GoStringer    */

func (n *BoxedInt) GoString() string {
    return fmt.Sprintf("%T{value:%s}", n, strconv.Itoa(n.value))
}

/*    implement fmt.Stringer    */

func (n *BoxedInt) String() string {
    return strconv.Itoa(n.GetInt())
}

/*    implement json.Marshaling    */

// MarshalJSON implements json.Marshaler iface.
// It encodes the boxed "int" value as JSON number.
func (n *BoxedInt) MarshalJSON() ([]byte, error) {
    return json.Marshal(n.GetInt())
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the boxed "int" value.
func (n *BoxedInt) UnmarshalJSON(b []byte) (err error) {
    var v int
    if err = json.Unmarshal(b, &v); err == nil {
        n.SetInt(v)
    }
    return
}




// An LockedInt is a BoxedInt with concurrent safe capabilities using a sync.Locker.
type LockedInt struct {
    sync.RWMutex    // <- embed
    BoxedInt        // <- embed
}

/*    Constructors    */

func NewLockedInt(value int) Inter {
    v := new(LockedInt)
    v.value = value    // attr provided by BoxedInt
    return v
}

/*    implement Inter    */

// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. SetInt(int)) like in other OOP languages happen.
func (n *LockedInt) SetInt(value int) {
    log.Println("LockedInt.Set")
    n.Lock()    // adquire-release write lock
    defer n.Unlock()
    // concurrent safe writer
// calling super embed
//    n.BoxedInt.SetInt(value)
// performance boosted
    n.value = value
}

// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. GetInt()) like in other OOP languages happen.
func (n *LockedInt) GetInt() int {
    log.Println("LockedInt.Get")
    n.RLock()    // adquire-release read lock
    defer n.RUnlock()
    // concurrent safe read
// calling super embed
//    return n.BoxedInt.GetInt()
// performance boosted
    return n.value
}

/*    implement Decrementer    */

// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. Dec()) like in other OOP languages happen.
func (n *LockedInt) Dec() Decrementer {
    log.Println("LockedInt.Dec")
    n.Lock()    // adquire-release write lock
    defer n.Unlock()
    // concurrent safe writter
// calling super embed
    n.BoxedInt.Dec()
// performance boosted
//    n.value--
    return n
}

/*    implement fmt.GoStringer    */

// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. GetInt()) like in other OOP languages happen.
func (n *LockedInt) GoString() string {
    // hide mutex
    return fmt.Sprintf("%T{value:%s}", n, strconv.Itoa(n.GetInt()))
}

/*    implement fmt.Stringer    */

// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. GetInt()) like in other OOP languages happen.
func (n *LockedInt) String() string {
    // hide mutex
    return strconv.Itoa(n.GetInt())
}

// MarshalJSON implements json.Marshaler iface.
// It encodes the boxed "int" value as JSON number.
// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. GetInt()) like in other OOP languages happen.
func (n *LockedInt) MarshalJSON() ([]byte, error) {
    return json.Marshal(n.GetInt())
}

// UnmarshalJSON implements json.Unmarshaler iface.
// It decodes JSON to the boxed "int" value.
// Overrides implementation provided by BoxedInt to be atomic.
// It's also necessary due to overrided methods here are not used by super embed structs (e.g. SetInt(int)) like in other OOP languages happen.
func (n *LockedInt) UnmarshalJSON(b []byte) (err error) {
    var v int
    if err = json.Unmarshal(b, &v); err == nil {
        n.SetInt(v)
    }
    return
}




// An LockedInter allows you to wrap any Inter instance to adquire concurrent safe capabilities using a sync.Locker.
type LockedInter struct {
    sync.RWMutex    // <- embed
    value    Inter  // wrapped instance
}

/*    Constructors    */

func NewLockedInter(value Inter) Inter {
    checkIntAtomicity(value)
    v := new(LockedInter)
    v.value = value
    return v
}

/*    implement Inter    */

func (n *LockedInter) SetInt(value int) {
    n.Lock()    // adquire-release write lock
    defer n.Unlock()
    // concurrent safe writer
    n.value.SetInt(value)
}

func (n *LockedInter) GetInt() int {
    n.RLock()    // adquire-release read lock
    defer n.RUnlock()
    // concurrent safe read
    return n.value.GetInt()
}

/*    implement fmt.GoStringer    */

func (n *LockedInter) GoString() string {
    // try first with type cast
    if gos, ok := n.value.(fmt.GoStringer); ok {
        n.RLock()    // adquire-release read lock
        defer n.RUnlock()
        return fmt.Sprintf("%T{value:%s}", n, gos.GoString())
    }
    // already safe in GetInt()
    return fmt.Sprintf("%T{value:%s}", n, strconv.Itoa(n.GetInt()))
}

/*    implement fmt.Stringer    */

func (n *LockedInter) String() string {
    // try first with type cast
    if s, ok := n.value.(fmt.Stringer); ok {
        n.RLock()    // adquire-release read lock
        defer n.RUnlock()
        return s.String()
    }
    // already safe in GetInt()
    return strconv.Itoa(n.GetInt())
}




func IsAtomicallySafe(value interface{}) bool {
    if value == nil {
        panic("atomic: nil AtomicInt detected.")
    }
    // safe type casting
    if _, ok := value.(sync.Locker); ok {
        return true
    }
    return false
}

func checkIntAtomicity(value Inter) {
    if IsAtomicallySafe(value) {
        panic("atomic: Is counterproductive to wrap one AtomicInt inside another.")
    }
}



func init() {
    log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
    log.SetOutput(os.Stdout)
}

func main(){
    var a Inter = &LockedInt{}
    fmt.Println(a) // 0
    printIs(a) // 3x true
    a.SetInt(9)
    fmt.Println(a) // 9
    if i, ok := a.(Incrementer); ok {
        i.Inc().Inc()
        fmt.Println(a) // 11
    }
    
    w := NewLockedInt(5)
    fmt.Println(w) // 5
    printIs(w)
    w.SetInt(3)
    fmt.Println(w) // 3
    if i, ok := w.(Decrementer); ok {
        i.Dec().Dec()
        fmt.Println(w) // 1
    }
    
    n := Int(5)
    fmt.Println(n) // 5
    
    nw := NewLockedInter(n)
    fmt.Println(nw) // 5
    printIs(nw)
    nw.SetInt(19)
    fmt.Println(nw) // 19
    fmt.Println(n) // 19
    
    NewLockedInter(nw)  // panic
}

func printIs(v interface{}) {
    fmt.Printf("Is %T Inter? %v\n", v, IsInter(v))
    fmt.Printf("Is %T Incrementer? %v\n", v, IsIncrementer(v))
    fmt.Printf("Is %T Decrementer? %v\n", v, IsDecrementer(v))
    var is bool
    _, is = v.(json.Marshaler)
    fmt.Printf("Is %T json.Marshaller? %v\n", v, is)
    _, is = v.(json.Unmarshaler)
    fmt.Printf("Is %T json.Unmarshaler? %v\n", v, is)
    fmt.Printf("fmt.GoStringer: %#v\n", v)
}
