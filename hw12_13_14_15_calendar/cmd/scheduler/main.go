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
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/rabbit"
	genstorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.toml", "Path to configuration file")
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

	var storage models.Storage
	storage, err = genstorage.NewStorage(ctx, log, config.Storage)
	if err != nil {
		log.Error("failed to create storage: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	addr := rabbit.BuildAMQPUrl(
		config.AMQP.Host, config.AMQP.Port, config.AMQP.Username, config.AMQP.Password,
	)
	server := rabbit.NewScheduler(config.Scheduler, storage, log)

	if err := server.Run(ctx, addr, config.AMQP.QueueName); err != nil {
		log.Error("failed to run scheduler: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
