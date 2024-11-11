package helpers

import (
	"reflect"

	"golang.org/x/exp/constraints"
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

func IsOfTypeInterface(inst interface{}, iface interface{}) bool {
	instType := reflect.TypeOf(inst)
	ifaceType := reflect.TypeOf(&iface).Elem() // .Elem() to get the interface type

	return instType.Implements(ifaceType)
}
