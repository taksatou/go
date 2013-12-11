# Functional

This is experimental package for functional programming with Go.
It provides a bunch of higher order functions.

example:

```
package main

import (
	"fmt"
	"github.com/taksatou/go/functional"
)

func addInt(a, b int) int {
	return a + b
}

func addFloat(a, b float64) float64 {
	return a + b
}

func main() {
	// declare a variable for curried function
	var curriedAdd func(int) func(int) int

	// curry add function and set it into curriedAdd
	functional.Curry(addInt, &curriedAdd)

	// apply 2 to curriedAdd and call the returned function with 1
	fmt.Println(curriedAdd(2)(1)) // => 3

	// Curry function works with any type of functions
	var curriedAddFloat func(float64) func(float64) float64
	functional.Curry(addFloat, &curriedAddFloat)
	fmt.Println(curriedAddFloat(0.5)(1.0)) // => 1.5
}
```
