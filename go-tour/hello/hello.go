package main

import (
	"fmt"

	"jpjofresm.com/greetings"
)

func main() {
	msg := greetings.Hello("Monito")
	fmt.Println(msg)
}
