package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	configFile, migrationVersion     string
	migrateUpgrade, migrateDowngrade bool
	ErrMigrationFailed               = errors.New("migration failed")
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.BoolVar(&migrateUpgrade, "migration-up", false, "perform migration update")
	flag.BoolVar(&migrateDowngrade, "migration-down", false, "perform migration downgrade")
	flag.StringVar(&migrationVersion, "migration-version", "", "set migration version")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("Can't parse config file, %v, exiting...", err)
	}
	log := logger.New(config.Logger.Level, config.Logger.Format)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	var storage app.Storage

	if config.Storage.Type == "memory" {
		storage = memorystorage.New()
	} else {
		sqlStorage, err := prepareSQLStorage(ctx, config.Storage.SQL, log)
		if err != nil {
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

	calendar := app.New(log, storage, fmt.Sprintf("%s:%v", config.HTTP.Host, config.HTTP.Port))
	server := internalhttp.NewServer(calendar)

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

func prepareSQLStorage(ctx context.Context, config SQLConf, log app.Logger) (*sqlstorage.Storage, error) {
	var command string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DBName,
	)
	if migrateUpgrade {
		command = "up"
	} else if migrateDowngrade {
		command = "down"
	}
	if err := performMigration(ctx, dsn, config.Driver, command, migrationVersion, log); err != nil {
		return nil, err
	}
	sqlStorage := sqlstorage.New()
	if err := sqlStorage.Connect(ctx, dsn, config.Driver); err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		return nil, err
	}
	return sqlStorage, nil
}

func performMigration(ctx context.Context, dsn, driver, command, version string, log app.Logger) error {
	if version != "" {
		command = fmt.Sprintf("%s-to", command)
	}
	db, err := goose.OpenDBWithDriver(driver, dsn)
	if err != nil {
		log.Info(fmt.Sprintf("goose: failed to open DB: %v\n", err))
		return ErrMigrationFailed
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error(fmt.Sprintf("goose: failed to close DB: %v\n", err))
		}
	}()

	if err := goose.RunContext(ctx, command, db, "./migrations"); err != nil {
		log.Error("goose %v: %v", command, err)
		return ErrMigrationFailed
	}
	return nil
}
