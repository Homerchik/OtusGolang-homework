package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger    LoggerConf `yaml:"logger"`
	HTTP      HTTPConf   `yaml:"http"`
	Storage   Storage    `yaml:"storage"`
	GRPC      GRPCConf   `yaml:"grpc"`
	AMQP      AMQPConf   `yaml:"amqp"`
	Scheduler Scheduler  `yaml:"scheduler"`
}

type LoggerConf struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCConf struct {
	Port int `yaml:"port"`
}

type Storage struct {
	Type string  `yaml:"type"`
	SQL  SQLConf `yaml:"sql"`
}

type SQLConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	Driver   string `yaml:"driver"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AMQPConf struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	QueueName string `yaml:"queueName"`
}

type Scheduler struct {
	MaxNotifyBefore int64 `yaml:"maxNotifyBefore"`
	ScanEvery       int64 `yaml:"scanEvery"`
	DeleteOlderThan int64 `yaml:"deleteOlderThan"`
	DeleteEvery     int64 `yaml:"deleteEvery"`
}

func NewConfig(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("decoding YAML: %w", err)
	}
	return config, nil
}
