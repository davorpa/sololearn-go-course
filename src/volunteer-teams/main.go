/*
Volunteer Teams +10 XP
Logical or boolean operators 3

Volunteers have been divided into 5 groups with equal number of members. If any group has fewer than 5 people, more volunteers must be recruited to fill the empty spots.
Write a program that takes the number of volunteers and outputs to the console how many volunteers need to be hired to have 5 equal groups.

Sample Input
24

Sample Output
1

Explanation
The nearest number to 24 that is multiple of 5 is 25, so we need 1 more volunteer (24+1=25) to get 5 equal groups.

Heads up!
Output 0 if we don't need additional volunteers.
Use the modulus (%) operator to get the division remainder.


@author davorpatech
@since  2021-09-15
*/

package main

import (
    "fmt"
    "os"
)

func main() {
    var n uint
    fmt.Print("Number of Volunteers? ")
    if _, err := fmt.Scanln(&n); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Println()
    
    fmt.Println(HiringReport(n))
}

func HiringReport(n uint) (s string) {
    volunteers := Volunteers(n)
    groups := volunteers.GetFormedGroups()
    rem := volunteers.GetRemainingPeople()
    switch {
        case groups == 0 && rem == 0:
            s = fmt.Sprintf("No volunteers yet. Try to recruit people!")
        case groups > 0 && rem == 0:
            s = fmt.Sprintf("You can form a total of %d %s.",
                groups, PluralizeUint(groups, "group", "groups"))
        case groups > 0 && rem > 0:
            s = fmt.Sprintf("You can form a total of %d %s and maybe make other if recruits %d %s more.",
                groups, PluralizeUint(groups, "group", "groups"),
                rem, PluralizeUint(rem, "volunteer", "volunteers"))
    }
    return
}

func PluralizeUint(value uint, singleLiteral, pluralLiteral string) string {
    if value != 1 {
        return pluralLiteral
    }
    return singleLiteral
}

const VOLUNTEERS_PER_GROUP uint = 5

type Volunteers uint

func (v Volunteers) GetFormedGroups() uint {
    return uint(v) / VOLUNTEERS_PER_GROUP
}

func (v Volunteers) GetRemainingPeople() uint {
    if mod := uint(v) % VOLUNTEERS_PER_GROUP; mod > 0 {
        return VOLUNTEERS_PER_GROUP - mod
    }
    return 0
}

func (v Volunteers) MatchesExactGroups() bool {
    return v.GetRemainingPeople() == 0
}
