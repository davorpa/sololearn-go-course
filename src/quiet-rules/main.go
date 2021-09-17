/*
Quiet rules
The if Statement +10 XP

Sundays in Switzerland are protected with Quiet Rules which make it illegal to pursue certain activities.
Taking the day of the week as input, output to the console "Obey the rules" if it is Sunday.

Sample Input
Sunday

Sample Output
Obey the rules

Hint:
Don't output anything if the input day isn't a Sunday.

@author davorpatech
@since  2021-09-17
*/

package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    var dow string
    fmt.Print("Input day of week: ")
    fmt.Scanln(&dow)
    fmt.Println()

    if matches, err := IsObeyRulesDayInSwitzerland(dow); err != nil {
        fmt.Fprintln(os.Stderr, "ERR!", err)
        os.Exit(1)
    } else if matches {
        fmt.Println("Obey the rules")
    }
}

var daysOfWeek = map[string]byte{
    "Monday": 0,
    "Tuesday": 0,
    "Wednesday": 0,
    "Thursday": 0,
    "Friday": 0,
    "Saturday": 0,
    "Sunday": 1,
    }

func IsObeyRulesDayInSwitzerland(dow string) (matches bool, err error) {
    found := false
    for k, v := range daysOfWeek {
        if equalsCaseInsensitive(k, dow) {
            matches = v == 1
            found = true
        }
    }
    if !found {
        err = fmt.Errorf("Invalid day of week: %s", dow)
    }
    return
}

func equalsCaseInsensitive(a, b string) bool {
    return strings.ToLower(a) == strings.ToLower(b)
}
