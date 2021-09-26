/*
Currency Converter. +50 XP

You are making a currency converter app.
Create a function called convert, which takes two parameters: the amount to convert, and the rate, and returns the resulting amount.
The code to take the parameters as input and call the function is already present in the Playground.
Create the function to make the code work.

Sample Input:
100
1.1

Sample Output:
110

Hint:
Converting 100 at the rate of 1.1 is equal to 100*1.1 = 110.

@author davorpatech
@since  2021-09-24
*/

package main

import (
    "errors"
    "fmt"
    "io"
    "math/big"
    "os"
    "strconv"
    "strings"
)

func main() {
    amount, err := AskF64("Input amount: ", 100)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }

    rate, err := AskF64("Input exchange rate: ", 1.1)
    if err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    }

    fmt.Printf("Exchange: %v",
        ConvertAmount(Amount(amount), FixedExchangeRate(rate)))
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
        return fallback, nil
    }

    v, err = strconv.ParseFloat(text, 64)
    if err != nil || v < 0 {
        err = fmt.Errorf("Invalid number: %q", text)
    }
    return
}

func ConvertAmount(amount Amount, rate RateExchanger) ExchangeResulter {
    return rate.ApplyRate(amount)
}

type ToBigFloater interface {
    BigFloat() *big.Float
}

func NewBigFloat() *big.Float {
    return new(big.Float).SetMode(
        big.ToNearestEven).SetPrec(256)
}

type Amount float64

func (a Amount) BigFloat() *big.Float {
    return NewBigFloat().SetFloat64(
        float64(a))
}

type RateExchanger interface {
    ApplyRate(amount Amount) ExchangeResulter
}

type FixedExchangeRate float64

func (er FixedExchangeRate) BigFloat() *big.Float {
    return NewBigFloat().SetFloat64(
        float64(er))
}

func (er FixedExchangeRate) ApplyRate(amount Amount) ExchangeResulter {
    exchange, _ := NewBigFloat().Mul(
        amount.BigFloat(), 
        er.BigFloat()).Float64()
    return exchangeResult{
        amount,
        Amount(exchange),
    }
}

type ExchangeResulter interface {
    GetSourceAmount() Amount
    GetExchangedAmount() Amount
    String() string
}

type exchangeResult struct {
    sourceAmount Amount
    exchangedAmount Amount
}

func (r exchangeResult) GetSourceAmount() Amount {
    return r.sourceAmount
}

func (r exchangeResult) GetExchangedAmount() Amount {
    return r.exchangedAmount
}

func (r exchangeResult) String() string {
    return fmt.Sprintf("%f => %f",
        r.sourceAmount,
        r.exchangedAmount)
}
