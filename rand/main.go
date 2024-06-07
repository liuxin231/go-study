package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// Int returns a non-negative pseudo-random int from the default
	randNum := rand.Int()
	fmt.Println("randNum: ", randNum, " len: ", len(fmt.Sprintf("%d", randNum)))

	// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n)
	randNum1 := rand.Int63()
	fmt.Println("randNum1: ", randNum1, " len: ", len(fmt.Sprintf("%d", randNum1)))
}
