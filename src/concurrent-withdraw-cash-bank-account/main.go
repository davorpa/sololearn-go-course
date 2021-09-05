/*
Concurrent Bank Account Withdrawing +50 XP

You are working on a banking app. 

The code declares a BankAccount struct with a balance field. 
You need to add a withdraw() method for the BankAccount struct. It should take an integer argument and decrease the balance of the Bank Account by the given amount.
In case there is not enough money in the account, the method should output "Insufficient Funds".

Take care to prevent more than one operation is not performed at same time.

Hint:
The code in main() takes the amount to withdraw as input, and calls the withdraw() method.

@author davorpatech
@since 2021-06-09
*/
package main
import (
    "errors"
    "fmt"
    "sync"
)

func main() {
    acc := NewBankAccount("James Smith", 100000)
    
    fmt.Printf("Hello, %v\n", acc.holder)
    var amount int
    fmt.Print("How much money do you want to withdraw? ")
    fmt.Scanln(&amount)
    
    if err := acc.Withdraw(amount); err != nil {
        switch werr := errors.Unwrap(err); werr {
            case InsufficientFundsErr:
                fmt.Println(werr)
            default:
                fmt.Println(err)
        }
    }
    fmt.Println(acc)
}

type BankAccount struct {
    // RWMutex: single-writer & multi-reader locks
    mux        sync.RWMutex
    holder     string
    balance    int
}

func NewBankAccount(holder string, balance int) *BankAccount {
    acc := new(BankAccount)
    acc.holder = holder
    acc.balance = balance
    return acc
}

func (acc BankAccount) String() string {
    // concurrent multi read protection
    acc.mux.RLock()
    defer acc.mux.RUnlock()
    // format
    return fmt.Sprint("{", acc.holder, " ", acc.balance, "}")
}

func (acc BankAccount) GetHolder() string {
    // concurrent multi read protection
    acc.mux.RLock()
    defer acc.mux.RUnlock()
    // safe read
    return acc.holder
}

func (acc BankAccount) GetBalance() int {
    // concurrent multi read protection
    acc.mux.RLock()
    defer acc.mux.RUnlock()
    // safe read
    return acc.balance
}

func (acc *BankAccount) Withdraw(amount int) (err error) {
    _, _, err = acc.operate(amount, "Withdraw", withdrawAmount)
    return
}

func (acc *BankAccount) Deposit(amount int) (err error) {
    _, _, err = acc.operate(amount, "Deposit", depositAmount)
    return 
}

// operate common pattern
func (acc *BankAccount) operate(amount int, action string, operation baccOperateFunc) (balance, newBalance int, err error) {
    // on exit, try transform into bussiness errors
    defer func() {
        if e := recover(); e != nil {
            // recover panics that their cast match
            if werr, ok := e.(bankAccountErr); ok {
                err = werr
            } else {
                panic(e) // unknown
            }
        }
        if err == nil {
            return
        }
        // wrap
        err = bankAccountErr{action, err}
    }()
    
    switch true {
        case amount < 0:
            err = NegativeAmountErr
            panic(err)
        case amount == 0:
            err = MissingAmountErr
        default:
            // concurrent single write protection
            acc.mux.Lock()
            defer acc.mux.Unlock()
            // safe operate
            balance, newBalance, err = operation(acc, amount)
    }
    return
}


/*      stateless & concurent unsafe helpers      */

type baccOperateFunc func(acc *BankAccount, amount int) (balance, newBalance int, err error)

func withdrawAmount(acc *BankAccount, amount int) (balance, newBalance int, err error) {
    if balance = acc.balance; balance < amount {
        err = InsufficientFundsErr
    } else {
        acc.balance -= amount
        newBalance, err = acc.balance, nil
    }
    return 
}

func depositAmount(acc *BankAccount, amount int) (balance, newBalance int, err error) {
    balance, err = acc.balance, nil
    acc.balance -= amount
    newBalance = acc.balance
    return
}



// =============================================
// Errors causes

var (
    NegativeAmountErr = errors.New("Negative Amount")
    MissingAmountErr = errors.New("Missing Amount")
    InsufficientFundsErr = errors.New("Insufficient Funds")
)



// =============================================
// Business Model Error

type bankAccountErr struct {
    action    string
    cause     error
}

func (e bankAccountErr) Action() string {
    return e.action
}

// Error implements "builtin.error" iface and provides the formatted string describing the cause and the action/operation that it was being done when the error happened.
func (e bankAccountErr) Error() string {
    return fmt.Sprint(e.cause, " to ", e.action)
}

// Unwrap implements "errors.Unwrapper" iface providing access to specific cause.
func (e bankAccountErr) Unwrap() error {
    return e.cause
}
