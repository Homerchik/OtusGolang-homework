package storage

import (
	"context"
	"fmt"

	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/config"
	"github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/homerchik/OtusGolang-homework/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(ctx context.Context, log app.Logger, storageConfig config.Storage) (models.Storage, error) {
	if storageConfig.Type == "memory" {
		return memorystorage.New(), nil
	}
	sqlConf := storageConfig.SQL
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		sqlConf.Host, sqlConf.Port, sqlConf.Username, sqlConf.Password, sqlConf.DBName,
	)
	sqlStorage := sqlstorage.New()
	if err := sqlStorage.Connect(ctx, dsn, sqlConf.Driver); err != nil {
		log.Error("failed to connect to storage: " + err.Error())
		return nil, err
	}
	return sqlStorage, nil
}
