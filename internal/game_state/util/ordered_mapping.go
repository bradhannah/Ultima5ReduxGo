package util

// this is inefficient, but handy so small sets of data that need a mapping
// to a string from a generic typed variable - but it is ordered.
// It's basically a stripped down `map` that orders stuff

type IdToString[T comparable] struct {
	Id           T
	FriendlyName string
}

type OrderedMapping[T comparable] []IdToString[T]

func (o *OrderedMapping[T]) GetIndex(thing T) int {
	for i, val := range *o {
		if val.Id == thing {
			return i
		}
	}
	return -1
}

func (o *OrderedMapping[T]) GetById(thing T) *IdToString[T] {
	for _, val := range *o {
		if val.Id == thing {
			return &val
		}
	}

	return nil
}
