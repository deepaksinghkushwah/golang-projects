package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	runtime.GOMAXPROCS(2)

	c := make(chan string)
	d := make(chan string)
	go makeCakeAndSend(c, "vanilla", 2)
	go makeCakeAndSend(d, "choco", 3)

	go receiveAndPack(c, d)
	//time.Sleep(2 * 1e9)
	fmt.Printf("Process took %s", time.Since(start))
	//fmt.Println(time.Now().Format("January 01, 2006 3:4:5 pm"))
}

func makeCakeAndSend(c chan string, flavour string, count int) {
	defer close(c)

	for i := 0; i <= count; i++ {
		cakeName := flavour + " cake " + strconv.Itoa(i)
		c <- cakeName
	}
}

func receiveAndPack(c chan string, d chan string) {
	cClosed, dClosed := false, false
	for {
		if cClosed && dClosed {
			return
		}
		fmt.Println("Waiting for a new cake ...")
		select {
		case cake, ok := <-c:
			if ok == false {
				cClosed = true
				fmt.Println(" ... vanila channel closed!")
			} else {
				fmt.Println(cake)
			}
		case cake, ok := <-d:
			if ok == false {
				dClosed = true
				fmt.Println(" ... choco channel closed!")
			} else {
				fmt.Println(cake)
			}
		default:
			fmt.Println(" ... all channels closed!")
		}
	}

}
