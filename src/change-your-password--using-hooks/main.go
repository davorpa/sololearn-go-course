/*
Change Your Password. Comparison Math Operators +10 XP

Reset account passwords should not be the same as the old password.

Task
Write a program that takes old and new passwords as inputs and outputs to the console the boolean value (true or false) that shows whether or not they match.

Sample Input
hl8792m
hl8792m

Sample Output
true

Hint:
Remember that all comparison operators always return true or false.
For example: console.log(5>4) outputs true.


@author davorpatech
@since  2021-09-12
*/

package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    var oldPass, newPass string
    fmt.Print("Your password: ")
    if _, err := fmt.Scanln(&oldPass); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    fmt.Print("Your new password: ")
    if _, err := fmt.Scanln(&newPass); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    isInvalidPass, setPass := useHistoryPasswordEffect(
        strings.Contains,
        oldPass)
    if isInvalidPass(newPass) {
        fmt.Println("Sorry! The password appears in your historial.")
    } else {
        setPass(newPass)
        fmt.Println("Password changed successfully.")
    }
}

func useHistoryPasswordEffect(
        comparator Comparator,
        historial... string,
    ) (
        Checker, 
        Setter,
        ) {
    // init memo
    memo := []string{}
    if len(historial) > 0 {
        memo = append(memo, historial...)
    }
    
    // hook accesors
    check := func(pwd string) (invalid bool) {
        for _, v := range memo {
            if comparator(v, pwd) {
                invalid = true
            }
        }
        return
    }
    set := func(pwd string) {
        memo = append(memo, pwd)
    }
    
    return check, set
}

type Comparator func (v, s string) bool
type Checker func(s string) bool
type Setter func(s string)
