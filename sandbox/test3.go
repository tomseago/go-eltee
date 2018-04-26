package main

import (
	"fmt"
	"github.com/nickysemenza/gola"
)

func main() {
	client := gola.New("localhost:9010")
	//defer client.Close()

	// # get DMX on universe 1
	data := make([]byte, 512)
	data[0] = 1
	data[1] = 2
	data[2] = 3

	for i := 0; i < 10; i++ {
		if b, err := client.SendDmx(1, data); err != nil {
			fmt.Printf("SendDmx: 1: %v", err)
		} else {
			fmt.Printf("SentDmx b=%v\n", b)
		}
	}
}
