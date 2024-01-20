package etc

import "fmt"

// inserts the given item at the specified index (negative indexes are from the end of the slice)
// note: reuses the caller's slice when possible (does not guarantee a copy)
// warn: panics if you are OOB
func InsertAt[T any](a []T, index int, value T) []T {

	// current length
	n := len(a)

	// negative indexes are from the end of the slice
	if index < 0 {
		index = n + index
	}

	// check for OOB
	if index < 0 || index > n {
		panic(fmt.Errorf("out of bounds: attempt to insert at %d for a %d length slice", index, n))
	}

	switch {
	case index == n:
		a = append(a, value)

	default:
		// extend our array by one element by duplicating the one before ourself
		a = append(a[:index+1], a[index:]...)
		a[index] = value
	}

	return a
}
