package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var (
		timeout time.Duration
		wg      sync.WaitGroup
	)
	pflag.DurationVar(&timeout, "timeout", time.Duration(10)*time.Second, "timeout for connection")
	pflag.Parse()
	args := pflag.Args()
	if len(args) != 2 {
		log.Fatalf("Usage: %s host port", os.Args[0])
	}
	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	wg.Add(2)
	go func() {
		for {
			if err := client.Receive(); err != nil {
				break
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			if err := client.Send(); err != nil {
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
