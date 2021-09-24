/*
Who Won The Match?
Multiple Parameters. +10 XP

You are given a program that takes Team 1 and Team 2 football teams goals as inputs accordingly.
Complete the function to take Team 1 and Team 2 goals as arguments and output the final result of the match:
- "Team 1 won", if Team 1's score is higher than Team 2's score
- "Team 2 won", if Team 2's score is higher than Team 1's score
- "Draw", if the scores are equal

Sample Input
3
4

Sample Output
Team 2 won


@author davorpatech
@since  2021-09-24
*/

package main

import (
    "errors"
    "fmt"
    "io"
    "math"
    "os"
    "strconv"
    "strings"
)

func main() {
    goalsTeam1, err := AskUint64("Goals Team 1: ", 3)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }
    goalsTeam2, err := AskUint64("Goals Team 2: ", 4)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }
    
    fmt.Println(finalResult(
        goalsTeam1, goalsTeam2))
}

func AskUint64(prompt string, fallback uint64) (v uint64, err error) {
    if len(prompt) > 0 {
        fmt.Print(prompt)
    }
    var text string
    _, err = fmt.Scan(&text)
    fmt.Println()
    if err != nil && !errors.Is(err, io.EOF) {
        return
    }

    if len(strings.TrimSpace(text)) == 0 {
        return fallback, nil
    }
    
    v, err = strconv.ParseUint(text, 10, 0)
    if err != nil {
        err = fmt.Errorf("Invalid number: %q", text)
    }
    return
}

func finalResult(goalsTeam1, goalsTeam2 uint64) string {
    var (
        diff int64 = int64(goalsTeam1) - int64(goalsTeam2)
        lbl  string
        args []interface{}
        sgoals = pluralizer("goal", "goals")
    )
    switch {
        case diff == 0:
            lbl = "Draw to %d %s"
            args = append(args, goalsTeam1, sgoals(goalsTeam1))
        case diff < 0:
            lbl = "Team 2 won by %d %s"
            args = append(args, int64(abs(diff)), sgoals(diff))
        case diff > 0:
            lbl = "Team 1 won by %d %s"
            args = append(args, int64(abs(diff)), sgoals(diff))
    }
    return fmt.Sprintf(lbl, args...)
}

func pluralizer(singleText, pluralText string) func(n interface{}) string {
    return func(n interface{}) string {
        if abs(n) == 1 {
            return singleText
        }
        return pluralText
    }
}

func abs(v interface{}) float64 {
    return math.Abs(ToFloat64(v))
}

type ToFloater64 interface{
    ToFloat64() float64
}

func ToFloat64(v interface{}) (n float64) {
    switch v.(type) {
        case ToFloater64: n = v.(ToFloater64).ToFloat64()
        case int: n = float64(v.(int))
        case *int: n = float64(*v.(*int))
        case int8: n = float64(v.(int8))
        case *int8: n = float64(*v.(*int8))
        case int16: n = float64(v.(int16))
        case *int16: n = float64(*v.(*int16))
        case int32: n = float64(v.(int32))
        case *int32: n = float64(*v.(*int32))
        case int64: n = float64(v.(int64))
        case *int64: n = float64(*v.(*int64))
        case uint: n = float64(v.(uint))
        case *uint: n = float64(*v.(*uint))
        case uint8: n = float64(v.(uint8))
        case *uint8: n = float64(*v.(*uint8))
        case uint16: n = float64(v.(uint16))
        case *uint16: n = float64(*v.(*uint16))
        case uint32: n = float64(v.(uint32))
        case *uint32: n = float64(*v.(*uint32))
        case uint64: n = float64(v.(uint64))
        case *uint64: n = float64(*v.(*uint64))
        case float32: n = float64(v.(float32))
        case *float32: n = float64(*v.(*float32))
        case float64: n = v.(float64)
        case *float64: n = *v.(*float64)
        default:
            panic(fmt.Sprintf("ToFloat64: type not suported: %T", v))
    }
    return
}
