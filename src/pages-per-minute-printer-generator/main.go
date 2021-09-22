/*
Pages per minute
The For loop. +10 XP

The printer prints 20 pages per minute.
Write a program that will output to the console how many pages will be printed each minute over the course of 6 minutes.

Output
20
40
60
80
100
120

Hint:
Don't forget to increment the counter (i++) for every iteration.


@author davorpatech
@since  2021-09-21
*/

package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    ctx := context.Background()
    minutes := uint(6)
    ppm, err := AskF64("Enter printer ppm: ", 20)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }

    printer := NewPrinter(ppm)
    for pages := range printer.Print(ctx, minutes) {
        fmt.Printf("Print progress: %4v pages\n", pages)
    }
}

func AskF64(prompt string, fallback float64) (v float64, err error) {
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
        v, err = fallback, nil
        return
    }

    v, err = strconv.ParseFloat(text, 64)
    if err != nil {
        err = fmt.Errorf("Invalid number: %q", text)
    }
    return
}

type printer struct {
    ppm    float64
}

func NewPrinter(ppm float64) *printer {
    if ppm < 0 {
        panic("Printer `ppm` must be positive number")
    }
    p := new(printer)
    p.ppm = ppm
    return p
}

func (p *printer) Print(ctx context.Context, minutes uint) chan float64 {
    if ctx == nil {
        panic("Print `context` is required")
    }
    ch := make(chan float64)
    refreshInterval := 100 * time.Millisecond
    ticker := time.NewTicker(refreshInterval)

    go func() {
        closing := false
        defer close(ch)
        defer ticker.Stop()

        var pages float64
        for i := uint(1); i <= minutes; i++ {
            pages = p.ppm * float64(i)
            select {
                // handle cancellation
                case <- ctx.Done():
                    closing = true
                    return
                // notify progress
                case <- ticker.C:
                    if ! closing {
                        ch <- pages
                    }
            }
        }
    }()

    return ch
}
