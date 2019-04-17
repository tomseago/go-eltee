package main

import (
	"fmt"
	"github.com/mgutz/ansi"
)

var titleColor = ansi.ColorFunc("blue")

func main() {
	fmt.Printf("%s\n", titleColor("ElTee Commander (ltc)"))

	lc := NewLocalContext()

	lc.Run()
}
