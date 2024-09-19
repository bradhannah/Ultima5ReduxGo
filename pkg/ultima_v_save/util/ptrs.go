package util

import "unsafe"

func MakeByteSliceFromUnsafePointer(ptr unsafe.Pointer, length int) []byte {
	// Ensure the length is non-negative
	if length < 0 {
		panic("length passed to makeByteSliceFromUnsafePointer is negative")
	}

	// Convert the unsafe.Pointer to a uintptr, then create a slice
	var slice []byte
	header := (*[1 << 30]byte)(ptr)[:length:length] // Create a slice pointing to the data
	slice = header
	return slice
}
