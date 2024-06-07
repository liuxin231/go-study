package main

import (
	"fmt"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Done()
	wg.Add(1)
	wg.Done()

	wg.Wait()
	fmt.Println("wait group wait end.")
}
