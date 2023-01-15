package config

import (
	"github.com/marki-eriker/event-listener/api"
	"github.com/marki-eriker/event-listener/db/postgres"
	"github.com/marki-eriker/event-listener/pkg/logger"
	"github.com/marki-eriker/event-listener/pkg/token"
	"github.com/spf13/viper"
)

// defaultPaths — пути по умолчанию для поиска файла конфигурации.
var defaultPaths = []string{".", "/usr/local/etc", "/etc", "/"}

// Load - загружает файл конфигурации
func Load(configName string) (*Config, error) {
	if configName == "" {
		return nil, ErrEmptyFileName
	}

	v, err := read(configName)
	if err != nil {
		return nil, err
	}

	c := fillConfig(v)

	return c, nil
}

func read(fileName string) (*viper.Viper, error) {
	opt := viper.New()

	for _, path := range defaultPaths {
		opt.AddConfigPath(path)
	}

	opt.SetConfigName(fileName)

	err := opt.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return opt, nil
}

func fillConfig(v *viper.Viper) *Config {
	c := &Config{
		Logger: &logger.Options{
			Level:       v.GetString("log.level"),
			ForceColors: v.GetBool("log.force_colors"),
		},
		JWT: &token.Options{
			Secret: v.GetString("token.secret"),
			Ttl:    v.GetDuration("token.ttl"),
		},
		Database: &postgres.Options{
			Addr:     v.GetString("db.addr"),
			User:     v.GetString("db.user"),
			Password: v.GetString("db.password"),
			Database: v.GetString("db.database"),
		},
		HTTPAPI: &api.HTTPAPIOptions{
			Addr:            v.GetString("api.http_api_addr"),
			Debug:           v.GetBool("api.debug"),
			ShutdownTimeout: v.GetDuration("api.shutdown_timeout"),
		},
		EventListener: &api.EventListenerOption{
			Addr: v.GetString("api.event_listener_addr"),
		},
	}

	return c
}
