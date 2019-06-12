package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	values := make(chan int, 3)
	defer close(values)
	go getRandomValues(values)
	go getRandomValues(values)
	fmt.Println(<-values)
	time.Sleep(1000 * time.Millisecond)
}

func getRandomValues(values chan int) {
	value := rand.Intn(10)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Random value: ", value)
	values <- value
	fmt.Println("Only Executes after another goroutine performs a receive on the channel")
}
