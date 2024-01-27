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

// does an in-situ reversal of any slice
func ReverseInSitu[S ~[]E, E any](s S) S {

	// get the length
	l := len(s)

	// only swap if at least 2 elements
	if l < 2 {
		return s
	}

	// we need our midpoint
	h := l / 2

	// set l to the last index, not the length
	l--

	// do an in-place swap of elements
	for i := 0; i < h; i++ {
		s[i], s[l-i] = s[l-i], s[i]
	}

	// purely to make composing these easier
	return s
}

// returns the reverse of the input slice (the output is separate from the input)
func Reverse[S ~[]E, E any](s S) S {
	return ReverseInSitu(Copy(s))
}

// creates a copy of the incoming slice, optionally with additional items appended
func Copy[S ~[]E, E any](source S, more ...E) S {
	c := make(S, 0, len(source)+len(more))
	c = append(c, source...) // subtle: my tests indicate that the built-in copy and append are identical in performance
	c = append(c, more...)
	return c
}

// removes the specified subset from the span (overwrites the slice memory)
func RemoveSpanInSitu[S ~[]E, E any](slice S, startIndex, count int) S {

	// check for invalid index or count
	if startIndex < 0 || count <= 0 || startIndex+count > len(slice) {
		panic("index or count out of bounds")
	}

	// use slicing to remove the subset
	return append(slice[:startIndex], slice[startIndex+count:]...)
}

// creates a new slice sans the given sub-span
func RemoveSpan[S ~[]E, E any](slice S, startIndex, count int) S {

	// check for invalid index or count
	if startIndex < 0 || count <= 0 || startIndex+count > len(slice) {
		panic("index or count out of bounds")
	}

	// use slicing to remove the subset
	return append(Copy(slice[:startIndex]), slice[startIndex+count:]...)
}

// inserts a series of copies of the specified element at index
func InsertSliceAt[S ~[]E, E any](slice S, subslice S, index int) S {

	// check for invalid index or count
	if index < 0 || index > len(slice) {
		panic("index out of bounds")
	}

	return append(slice[:index], append(subslice, slice[index:]...)...)
}

// returns a slice that is a simple run of the given element
func MakeRun[S ~[]E, E any](element E, count int) (run S) {

	// check for invalid index or count
	if count < 0 {
		panic("index or count out of bounds")
	}

	run = make(S, 0, count)
	for i := count; i != 0; i-- {
		run = append(run, element)
	}

	return
}

// inserts a series of copies of the specified element at index
func InsertRunAt[S ~[]E, E any](slice S, index, count int) S {

	// check for invalid index or count
	if index < 0 || index > len(slice) {
		panic("index or count out of bounds")
	}

	// insert the duplicated run at the specified index
	InsertSliceAt(slice, MakeRun[S](slice[index], count), index)

	return slice
}
