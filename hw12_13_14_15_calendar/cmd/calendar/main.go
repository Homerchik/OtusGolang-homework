package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/homerchik/hw12_13_14_15_calendar/internal/app"
	"github.com/homerchik/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/homerchik/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/homerchik/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/homerchik/hw12_13_14_15_calendar/internal/storage/sql"
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

	config := NewConfig()
	log := logger.New(config.Logger.Level)
	ctx, cancel := signal.NotifyContext(context.Background(),
	syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage app.Storage
	if config.Storage.Type == "memory" {
		storage = memorystorage.New()
	} else {
		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
			config.Storage.SQL.Host, config.Storage.SQL.Port, config.Storage.SQL.Username, config.Storage.SQL.Password, config.Storage.SQL.DbName,
		)
		sqlStorage := sqlstorage.New()
		if err := sqlStorage.Connect(ctx, connStr, config.Storage.SQL.Driver); err != nil {
			log.Error("failed to connect to storage: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
		defer func() {
			if err := sqlStorage.Close(ctx); err != nil {
				log.Error("failed to close storage: " + err.Error())
			}
		}()
		storage = sqlStorage
	}
	calendar := app.New(log, storage)

	server := internalhttp.NewServer(log, calendar)

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
		os.Exit(1) //nolint:gocritic
	}
}
