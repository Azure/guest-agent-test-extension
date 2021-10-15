package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("GATESTEXT SERVICE STARTED")

	for true {
		time.Sleep(1 * time.Second)
		fmt.Printf("GATESTEXT SERVICE RUNNING")
	}
}