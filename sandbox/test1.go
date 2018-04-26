package main

import (
	"fmt"
	"github.com/nickysemenza/gola"
)

func main() {
	// client := gola.New("localhost:9010")
	client := gola.New("10.0.1.121:9010")
	defer client.Close()

	// # get DMX on universe 1
	if x, err := client.GetDmx(1); err != nil {
		fmt.Printf("GetDmx: 1: %v", err)
	} else {
		fmt.Printf("GetDmx: 1: %v", x.Data)
	}
}
