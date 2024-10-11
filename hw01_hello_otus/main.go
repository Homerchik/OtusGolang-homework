package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	// Place your code here.
	greet := "Hello, OTUS!"
	fmt.Print(reverse.String(greet))
}
