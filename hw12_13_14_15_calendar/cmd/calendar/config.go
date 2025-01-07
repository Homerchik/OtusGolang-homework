package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Http HttpConf
	Storage Storage
}

type LoggerConf struct {
	Level string
	// TODO
}

type HttpConf struct {
	Host string
	Port int
}

type Storage struct {
	Type string
	SQL SQLConf
}

type SQLConf struct {
	Host string
	Port int
	Dialect string
	Username string
	Password string
}

func NewConfig() Config {
	return Config{}
}