/*
BOUNCY NUMBER CHALLENGE

A bouncy number is any non-negative integer that is neither increasing nor decreasing. 
Ex-1234, not a bouncy number(ascending order)
Ex-4321, not a bouncy number(descending order)
Ex-7264, bouncy number

@author davorpatech
@since 2021-04-31
*/

package main

import (
    . "fmt"
    "os"
    "sort"
)


func main() {
    //
    // try to read input as number
    //
    Print("Enter a number to check: ")
    var num int64
    if _, err := Scanln(&num); err != nil {
        //panic(err)
        Fprintln(os.Stderr, "⚠", err)
        os.Exit(1)
        return
    }
    
    //
    // process
    //
    status, err := IsBouncyInt64(num)
    if err != nil {
        //panic(err)
        Fprintln(os.Stderr, "⚠", err)
        os.Exit(1)
        return
    }
    if status == Bouncy {
        Println(num, "IS a bouncy number.")
    } else if status == NotBouncyAsc {
        Println(num, "NOT a bouncy number (increasing order).")
    } else if status == NotBouncyDesc {
        Println(num, "NOT a bouncy number (decreasing order).")
    }
}





//   A bouncy number is any non-negative integer that its digits is neither increasing nor decreasing.
func IsBouncyInt64(num int64) (BouncyStatus, error) {
    // preconditions
    if num < 0 {
        return Bouncy, Errorf("bouncy: Uncheckable number: %v", num)
    }
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text), nil
}

func IsBouncyInt32(num int32) (BouncyStatus, error) {
    // preconditions
    if num < 0 {
        return Bouncy, Errorf("bouncy: Uncheckable number: %v", num)
    }
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text), nil
}

func IsBouncyInt(num int) (BouncyStatus, error) {
    // preconditions
    if num < 0 {
        return Bouncy, Errorf("bouncy: Uncheckable number: %v", num)
    }
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text), nil
}

func IsBouncyUInt(num uint) BouncyStatus {
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text)
}

func IsBouncyUInt32(num uint32) BouncyStatus {
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text)
}

func IsBouncyUInt64(num uint64) BouncyStatus {
    // to string
    text := Sprint(num)
    // check
    return IsBouncyDigital(text)
}


//   BouncyStatus enum
type BouncyStatus int
const (
    NotBouncyDesc BouncyStatus = iota // 0
    Bouncy                       // 1
    NotBouncyAsc                 // 2
)
// String implements fmt.Stringer iface.
func (status BouncyStatus) String() (s string) {
    switch status {
        case NotBouncyDesc: s = "Not Bouncy Number (decreasing order)"
        case Bouncy: s = "Bouncy Number"
        case NotBouncyAsc: s = "Not Bouncy Number (increasing order)"
    }
    return
}
// GoString implements fmt.GoStringer iface in charge to format value when %#v verb is used.
func (status BouncyStatus) GoString() (s string) {
    switch status {
        case NotBouncyDesc: s = "%T<%d-NotBouncyDesc>"
        case Bouncy: s = "%T<%d-Bouncy>"
        case NotBouncyAsc: s = "%T<%d-NotBouncyAsc>"
    }
    return Sprintf(s, status, status)
}


//   IsBouncyDigital checks byte to byte if not match
func IsBouncyDigital(s string) BouncyStatus {
    // to byte slice
    bytes := ByteSlice([]byte(s))
    // check ascending sort
    if sort.IsSorted(bytes) {
        return NotBouncyAsc
    }
    // check descending sort
    sort.Stable(sort.Reverse(bytes))
    if string(bytes) == s {
        return NotBouncyDesc
    }
    return Bouncy
}



//   ByteSlice alias
type ByteSlice []byte

// Implements Sort Interface, sorting in increasing order.
func (this ByteSlice) Len() int {
    return len(this)
}
func (this ByteSlice) Less(i, j int) bool {
    return this[i] < this[j]
}
func (this ByteSlice) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}
