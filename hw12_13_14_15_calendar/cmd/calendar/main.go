package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/config"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	internalhttp "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/server/http"
	genstorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
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

	addr := fmt.Sprintf("%s:%v", config.HTTP.Host, config.HTTP.Port)
	server := internalhttp.NewServer(addr, log, storage)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
