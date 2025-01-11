package main

import "os"
import "fmt"
import "gopkg.in/yaml.v3"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	Http HttpConf     `yaml:"http"` 		
	Storage Storage   `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

type HttpConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Storage struct {
	Type string  `yaml:"type"`
	SQL SQLConf  `yaml:"sql"`
}

type SQLConf struct {
	Host string      `yaml:"host"`
	Port int         `yaml:"port"`
	DbName string    `yaml:"db_name"`
	Driver string    `yaml:"driver"`
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
}

func NewConfig(filepath string) (Config, error) {
	file, err := os.Open(filepath)
    if err != nil {
        return Config{}, fmt.Errorf("opening file: %v", err)
    }
    defer file.Close()

    decoder := yaml.NewDecoder(file)
    var config Config
    if err := decoder.Decode(&config); err != nil {
        return Config{}, fmt.Errorf("decoding YAML: %v", err)
    }
	return config, nil
}