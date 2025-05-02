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

func IsOfTypeInterface(inst interface{}, iface interface{}) bool {
	instType := reflect.TypeOf(inst)
	ifaceType := reflect.TypeOf(&iface).Elem() // .Elem() to get the interface type

	return instType.Implements(ifaceType)
}

func OneInXOdds(odds int) bool {
	//rand.Seed(uint64(time.Now().UnixNano()))

	return rand.Intn(odds) == 0
}

func RandomIntInRange(min, max int) int {
	if min > max {
		panic("min cannot be greater than max")
	}
	//rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator
	return rand.Intn(max-min+1) + min
}

func PickOneOf[T any](a, b T) T {
	if OneInXOdds(2) {
		return a
	}
	return b
}

// func Swap[T comparable](a, b T) (T, T) {
// 	return b, a
// }
