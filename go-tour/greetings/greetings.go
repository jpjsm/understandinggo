package greetings

import "fmt"

func Hello(name string) string {
	msg := fmt.Sprintf("Hi %v, welcome!", name)
	return msg
}
