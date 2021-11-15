package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("GATESTEXT SERVICE STARTED")

	for true {
		time.Sleep(5 * time.Minute)
		fmt.Printf("GATESTEXT SERVICE RUNNING")
	}
}