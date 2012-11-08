package nimble

// Utility functions

import (
	"flag"
	"os"
	"strconv"
)

// Open file for writing, panic or error.
func OpenFile(fname string) *os.File {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	Fatal(err)
	return f
}

// Product of elements.
func Prod(size [3]int) int {
	return size[0] * size[1] * size[2]
}

// Wraps an index to [0, max] by adding/subtracting a multiple of max.
func Wrap(number, max int) int {
	for number < 0 {
		number += max
	}
	for number >= max {
		number -= max
	}
	return number
}

// Panics if a != b
func CheckEqualSize(a, b [3]int) {
	if a != b {
		Panic("Size mismatch:", a, "!=", b)
	}
}

func CheckUnits(a, b string) {
	if a != b {
		Panicf(`Unit mismatch: "%v" != "%v"`, a, b)
	}
}

// IntArg returns the idx-th command line as an integer.
func IntArg(idx int) int {
	val, err := strconv.Atoi(flag.Arg(idx))
	Fatal(err)
	return val
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}