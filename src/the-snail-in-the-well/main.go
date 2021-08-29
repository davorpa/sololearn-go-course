/*
The Snail in the Well. +10 XP


The snail climbs up 7 feet each day and slips back 2 feet each night.
How many days will it take the snail to get out of a well with the given depth?

Sample Input:
31

Sample Output:
6

Explanation:
Let's break down the distance the snail covers each day:
Day 1: 7-2=5
Day 2: 5+7-2=10
Day 3: 10+7-2=15
Day 4: 15+7-2=20
Day 5: 20+7-2=25
Day 6: 25+7=32
So, on Day 6 the snail will reach 32 feet and get out of the well at day, without slipping back that night.

Hint:
You can use a loop to calculate the distance the snail covers each day, and break the loop when it reaches the desired distance.

@author davorpatech
@since 2021-08-30
*/
package main

import (
	."fmt"
	"os"
	)

func main() {
	var depth uint
	Print("Input depth: ")
	if _, err := Scan(&depth); err != nil {
		Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	Println()
	
	days := CalculateDays(depth)
	Printf("The snail needs %d %s to get out of a well with a given depth of %d feet.\n", days, PluralizeUint(days, "day", "days"), depth)
}



func CalculateDays(depth uint) (days uint) {
	// we use an infinite loop because
	// stop condition happens between
	// 2 computational sentences
	for climb := 0; ; {
		// past days
		days++
		// day feet advance
		climb += SnailClimb(false)
		// depth reached ??
		if climb >= int(depth) {
			break
		}
		// night feet slipping
		climb += SnailClimb(true)
	}
	return
}


const (
	ADVANCE_FEET  = 7
	SLIPPING_FEET = 2
	)

func SnailClimb(slipping bool) int {
	if slipping {
		return -SLIPPING_FEET
	}
	return ADVANCE_FEET
}


func PluralizeUint(value uint, singleLiteral, pluralLiteral string) string {
	if value != 0 {
		return pluralLiteral
	}
	return singleLiteral
}
