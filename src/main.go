/*
Math Operators 1: Speed of Sound. +10 XP

Science fact! The speed of sound in air is about 343 meters per second.

Create a program that calculates and outputs to the console the distance the sound covers in 600 seconds.


@author davorpatech
@since 2021-08-30
*/
package main

import (
	."fmt"
	"os"
)

func main() {
	var seconds uint = 600
	Print("Input seconds: ")
	if _, err := Scan(&seconds); err != nil {
		Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	Println()
	
	distance := ComputeSoundDistanceByTime(seconds)
	Printf("The distance the sound covers in %d %s is: %d %s.\n",
		seconds, PluralizeUint(seconds, "second", "seconds"),
		distance, PluralizeUint(distance, "meter", "meters"))
}

const (
	MPS_SPEED_OF_SOUND_IN_AIR = 343
)

func ComputeSoundDistanceByTime(seconds uint) uint {
	return seconds * MPS_SPEED_OF_SOUND_IN_AIR
}

func PluralizeUint(value uint, singleLiteral, pluralLiteral string) string {
	if value != 0 {
		return pluralLiteral
	}
	return singleLiteral
}
