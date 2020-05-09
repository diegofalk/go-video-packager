package main

import (
	"fmt"
	"time"
)

func packagerRun() {
	for {
		id := <-data
		fmt.Printf("i=%v\n", id)
		time.Sleep(10 * time.Second)
	}
}
