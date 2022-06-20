package main

import (
	"fmt"
	"math"
)

func divide(x, y int) int {
	var sign bool
	var q, result int

	if x < 0 {
		sign = true
	} else {
		sign = false
	}

	x = int(math.Abs(float64(x)))

	if (y > x) || (y < 0) || (y == 0) {
		return 0
	}

	q = divide(x, 2*y)

	if (x - q*2*y) < y {
		result = 2 * q
	} else {
		result = 2*q + 1
	}

	if sign {
		return -result
	} else {
		return result
	}
}

func main() {
	fmt.Println(divide(-30, 6))
	fmt.Println(divide(-18000, 6))
	fmt.Println(divide(18, -1))
	fmt.Println(divide(18, 0))
}
