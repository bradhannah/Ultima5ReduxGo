package helpers

import (
	"reflect"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/rand"
)

// Max Generic Max function
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min Generic Min function
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// func Abs[T constraints.Integer | constraints.Float](a, b T) T {}

func AbsInt(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

func IsOfTypeInterface(inst, iface interface{}) bool {
	instType := reflect.TypeOf(inst)
	ifaceType := reflect.TypeOf(&iface).Elem() // .Elem() to get the interface type

	return instType.Implements(ifaceType)
}

func OneInXOdds(odds int) bool {
	return rand.Intn(odds) == 0
}

// HappenedByPercentLikely returns true if the event has happened by the given likelihood.
// 0% likelihood means it will never happen.
// 100% likelihood means it will always happen.
func HappenedByPercentLikely(likelihoodToSucceedPercent int) bool {
	if likelihoodToSucceedPercent >= 100 {
		return true
	}
	if likelihoodToSucceedPercent <= 0 {
		return false
	}
	return RandomIntInRange(0, 100) < likelihoodToSucceedPercent
}

func RandomIntInRange(min, max int) int {
	if min > max {
		panic("min cannot be greater than max")
	}
	return rand.Intn(max-min+1) + min
}

func PickOneOf[T any](a, b T) T {
	if OneInXOdds(2) {
		return a
	}
	return b
}

func FilterFromSlice[T any](s []T, keep func(T) bool) []T {
	nIndex := 0 // next write position

	for _, v := range s { // read position
		if keep(v) {
			s[nIndex] = v // overwrite; OK even when nIndex==read index
			nIndex++
		}
	}

	return s[:nIndex] // truncate tail
}

func RoundUp(f float32) int {
	return int(f + 0.5)
}
