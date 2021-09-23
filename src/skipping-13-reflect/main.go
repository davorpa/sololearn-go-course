/*
Skipping 13
Break and Continue. +10 XP

Many tall buildings, including hotels, skip the number 13 when numbering floors -- often going from floor 12 to floor 14. It is thought that the number 13 is unlucky.
Write a program that will number 15 rooms starting from 1, skipping the number 13. Output to the console each room number in separate line.

Hint:
Be attentive - considering the missing number, the number of last room should be greater than the count of rooms by 1.


@author davorpatech
@since  2021-09-22
*/

package main

import (
    "fmt"
    "reflect"
    "strings"
)

const countOfRooms int = 15
var unluckys = [...]int{13}

func main() {
    for n, roomNumber := 0, 1; n < countOfRooms; roomNumber++ {
        if ok, _ := Contains(unluckys, roomNumber); ok {
            continue
        }
        n++
        fmt.Println(roomNumber)
    }
}

//======================================
// UTILS. reflect
//

//
// `Contains` inspect using reflection if some elem is contained in any iterable data struct like slice, arrays, maps, strings...
//
func Contains(c interface{}, elem interface{}) (bool, int) {
    if c == elem {
        return true, 0
    }
    if c == nil && elem != nil {
        return false, -1
    }

    if kind := reflect.TypeOf(c).Kind(); kind == reflect.Slice || kind == reflect.Array {
        return ContainsSliced(c, elem)
    } else if kind == reflect.Map {
        return ContainsMapKey(c, elem)
    }

    // safe casting as strings
    if str, ok := c.(string); ok {
        if search, ok := elem.(string); ok {
            i := strings.Index(str, search)
            return i != -1, i
        }
    }

    return false, -1 // not found
}

//
// `ContainsSliced` inspect using reflection if some item is contained in any iterable data struct like slice or arrays. It panics for other types.
//
func ContainsSliced(c interface{}, elem interface{}) (found bool, pos int) {
    cval := reflect.ValueOf(c)
    if k := cval.Kind(); k != reflect.Slice && k != reflect.Array {
        panic("ContainsSliced only can process slices or static arrays")
    }
    found, pos = false, -1 // not found

    if c == nil && elem != nil {
        return
    }
    
    // iterate over array elements
    for i, l := 0, cval.Len(); i < l && !found; i++ {
        // XXX - maybe panics if sliced element points to an unexported struct field
        if checkContainsIdentity(cval.Index(i), elem) {
            found, pos = true, i
        }
    }
    return
}

//
// `ContainsMapKey` inspect using reflection if some item is contained as key by any map. It panics for other types.
//
func ContainsMapKey(m interface{}, key interface{}) (found bool, pos int) {
    return containsMap(m, key,
        mapIterEntryKey)
}

//
// `ContainsMapValue` inspect using reflection if some item is contained as value by any map. It panics for other types.
//
func ContainsMapValue(m interface{}, value interface{}) (found bool, pos int) {
    return containsMap(m, value,
        mapIterEntryValue)
}

type MapIterFieldFunc func(*reflect.MapIter) reflect.Value

func mapIterEntryKey(it *reflect.MapIter) reflect.Value {
    return it.Key()
}

func mapIterEntryValue(it *reflect.MapIter) reflect.Value {
    return it.Value()
}

func containsMap(m interface{}, item interface{}, fieldFun MapIterFieldFunc) (found bool, pos int) {
    mval := reflect.ValueOf(m)
    if mval.Kind() != reflect.Map {
        panic("ContainsMap only can process maps")
    }
    found, pos = false, -1 // not found

    if m == nil && item != nil {
        return
    }

    // iterate over map entries
    // is not ensured that position `i` be equals over same calls.
    for i, it := 0, mval.MapRange(); it.Next() && !found; i++ {
        // XXX - maybe panics if map entry field value points to an unexported struct field
        if checkContainsIdentity(fieldFun(it), item) {
            found, pos = true, i
        }
    }
    return
}

func checkContainsIdentity(value reflect.Value, item interface{}) (ok bool) {
    // XXX - panics if value points to an unexported struct field
    // see https://golang.org/pkg/reflect/#Value.Interface
    // XXX - use Value.CanInterface() ??
    if value.CanInterface() && value.Interface() == item {
        ok = true
    }
    return
}
