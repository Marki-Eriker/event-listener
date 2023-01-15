package config

import (
	c "github.com/kaatinga/const-errs"
	"github.com/marki-eriker/event-listener/api"
	"github.com/marki-eriker/event-listener/db/postgres"
	"github.com/marki-eriker/event-listener/pkg/logger"
	"github.com/marki-eriker/event-listener/pkg/token"
)

const ErrEmptyFileName c.Error = "empty config file name"

type Config struct {
	Logger        *logger.Options
	JWT           *token.Options
	Database      *postgres.Options
	HTTPAPI       *api.HTTPAPIOptions
	EventListener *api.EventListenerOption
}
