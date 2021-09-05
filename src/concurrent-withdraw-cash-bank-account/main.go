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
import "fmt"

type BankAccount struct {
  holder string
  balance int
}

func main() {
  acc := BankAccount{"James Smith", 100000}
  
  var amount int
  fmt.Scanln(&amount)
  
  acc.withdraw(amount)
  fmt.Println(acc)
}

