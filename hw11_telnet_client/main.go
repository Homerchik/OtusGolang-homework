package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Duration(10)*time.Second, "timeout for connection")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatalf("Usage: %s host port", os.Args[0])
	}
	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for {
			if err := client.Send(); err != nil {
				break
			}
		}
		stop()
	}()
	go client.Receive()
	<-ctx.Done()
}
