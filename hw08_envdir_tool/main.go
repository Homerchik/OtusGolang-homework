package main

import (
	"log"
	"os"
)

func main() {
	// Place your code here.
	args := os.Args
	if len(args) < 3 {
		log.Fatalf("No enough args have been passed. Passed args quantity should at least 2")
	}
	if env, err := ReadDir(args[1]); err == nil {
		code := RunCmd(args[2:], env)
		log.Printf("Command has been executed, return code is %v", code)
	} else {
		log.Fatal(err)
	}
}
