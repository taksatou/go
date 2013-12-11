package main

import (
	"fmt"
	"github.com/taksatou/go/functional"
)

func curryExample1() {
	// declare a variable for curried function
	var curriedAdd func(int) func(int) int

	// curry add function and set it into curriedAdd
	functional.Curry(add, &curriedAdd)

	// apply 2 to curriedAdd and call the returned function with 1
	fmt.Println(curriedAdd(2)(1)) // => 3
}

func curryExample2() {
	// Curry function works with any type of functions
	var curriedFold func(interface{}) func(interface{}, interface{}) (interface{}, error)
	functional.Curry(functional.Fold, &curriedFold)
	sumInt := curriedFold(func(a, b int) int { return a + b })
	fmt.Println(sumInt([]int{1, 2, 3}, nil)) // => 6

	sumFloat64 := curriedFold(func(a, b float64) float64 { return a + b })
	// initial value should be the same type as slice elements.
	// convert `1` to float64 value explicitly, since number literal is interpreted as integer
	fmt.Println(sumFloat64([]float64{1, 2, 3}, float64(1))) // => 6.1
}

func volume(a, b, c float64) float64 {
	return a * b * c
}

func add(a, b int) int {
	return a + b
}

func main() {
	curryExample1()
	curryExample2()
}
