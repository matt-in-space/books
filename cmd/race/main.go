package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	counter := 1

	go func() {
		mu.Lock()
		defer mu.Unlock()
		if counter > 0 {
			counter--
			fmt.Println("A won")
		}
	}()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		if counter > 0 {
			counter--
			fmt.Println("B won")
		}
	}()

	time.Sleep(1 * time.Second)
}

func doA() {
	for i := range 100 {
		fmt.Printf("Hello from A: %d\n", i)
		time.Sleep(10 * time.Millisecond)
	}
}

func doB() {
	for i := range 100 {
		fmt.Printf("Hello from B: %d\n", i)
		time.Sleep(10 * time.Millisecond)
	}
}
