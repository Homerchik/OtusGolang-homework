package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/config"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Can't parse config file, %v, exiting...", err)
	}
	log := logger.New(config.Logger.Level, config.Logger.Format)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	addr := rabbit.BuildAMQPUrl(
		config.AMQP.Host, config.AMQP.Port, config.AMQP.Username, config.AMQP.Password,
	)
	server := rabbit.NewSender(log)

	if err := server.Run(ctx, addr, config.Sender.RcvQueue, config.Sender.PushQueue); err != nil {
		log.Error("failed to run sender: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
