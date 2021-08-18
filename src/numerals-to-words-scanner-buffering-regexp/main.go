/*
NO NUMERALS +10 XP

You write a phrase and include a lot of number characters (0-9), but you decide that for numbers 10 and under you would rather write the word out instead. Can you go in and edit your phrase to write out the name of each number instead of using the numeral?

Task:
Take a phrase and replace any instances of an integer from 0-10 and replace it with the English word that corresponds to that integer.

Input Format:
A string of the phrase in its original form (lowercase).

Output Format:
A string of the updated phrase that has changed the numerals to words.

Sample Input:
    i need 2 pumpkins and 3 apples

Sample Output:
    i need two pumpkins and three apples


@author davorpatech
@since  2021-08-18
*/
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    "sync"
    )

func init() {
    log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
    log.SetOutput(os.Stdout) //os.Stderr
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanLines)
    
    fmt.Print("Enter your phase: ")
    var text strings.Builder
    for i := 0; scanner.Scan(); i++ {
        line := scanner.Text()
        // rebuild previous newline
        if i > 0 {
            text.WriteString("\n")
        }
        text.WriteString(line)
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }
    
    fmt.Println()
    fmt.Println("No numerals result:")
    fmt.Println(
        ReplaceNumerals(text.String()))
}


func ReplaceNumerals(text string) string {
    // lazy init regex compile
    enclosedDigitsRegexOnce.Do(
        initEnclosedDigitsRegex)
    // replace each match using a mapper
    return enclosedDigitsRegex.ReplaceAllStringFunc(
        text, numeralReplacer)
}


var enclosedDigitsRegexOnce sync.Once
var enclosedDigitsRegex *regexp.Regexp

func initEnclosedDigitsRegex() {
    enclosedDigitsRegex = regexp.MustCompile("(^|(?:\\b+))(\\d+)((?:\\b+)|$)")
}

func numeralReplacer(match string) string {
    switch match {
        case  "0": match = "zero";
        case  "1": match = "one";
        case  "2": match = "two";
        case  "3": match = "three";
        case  "4": match = "four";
        case  "5": match = "five";
        case  "6": match = "six";
        case  "7": match = "seven";
        case  "8": match = "eight";
        case  "9": match = "nine";
        case "10": match = "ten";
    }
    //don't touch if not match
    return match
}
