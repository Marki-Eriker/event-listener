package main

import (
	"flag"
	"github.com/marki-eriker/event-listener/api"
	"github.com/marki-eriker/event-listener/config"
	"github.com/marki-eriker/event-listener/db/postgres"
	"github.com/marki-eriker/event-listener/metric"
	"github.com/marki-eriker/event-listener/pkg/app"
	"github.com/marki-eriker/event-listener/pkg/encoder"
	"github.com/marki-eriker/event-listener/pkg/hasher"
	"github.com/marki-eriker/event-listener/pkg/logger"
	"github.com/marki-eriker/event-listener/pkg/recovery"
	"github.com/marki-eriker/event-listener/pkg/request"
	"github.com/marki-eriker/event-listener/pkg/token"
	"github.com/marki-eriker/event-listener/processor"
	"github.com/sirupsen/logrus"
	"net/http"
)

var configName string

func init() {
	flag.StringVar(&configName, "config", "service-dev", "configuration file name")
}

func main() {
	flag.Parse()

	opt, err := config.Load(configName)
	if err != nil {
		logrus.WithError(err).Fatal("unable to load service configuration")
	}

	m := metric.NewMock()

	log, err := logger.New(opt.Logger)
	if err != nil {
		logrus.WithError(err).Fatal("unable to configure logger")
	}

	db, err := postgres.New(opt.Database, m.AddDBQuery)
	if err != nil {
		log.WithError(err).Fatal("unable to configure DB")
	}

	request.Setup(log, m.AddResponse)
	recovery.Setup(m.AddPanic)

	s := &processor.Store{
		User:  db.User,
		Event: db.Event,
	}

	t := token.NewJWT(opt.JWT)
	h := hasher.NewBcryptHasher(0)
	e := encoder.NewMock()

	p := processor.New(log, s, t, h, e)

	httpSrv := api.NewHTTPAPI(p, log, opt.HTTPAPI, m, t)

	go func(srv *api.Server) {
		defer recovery.Restart("http_api")

		if err := srv.Serve(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatalf("Unable to serve HTTP API")
		}
	}(httpSrv)

	eventListener := api.NewEventListener(opt.EventListener, log, m, httpSrv.HandleEvent)

	go func(lis *api.Listener) {
		defer recovery.Restart("event_listener")

		if err := lis.Start(); err != nil {
			log.WithError(err).Fatalf("unable to start listen events")
		}
	}(eventListener)

	app.New(log).Run()

	eventListener.Stop()
	httpSrv.Stop()
	db.Stop()

	log.Warn("shutdown is executed")
}
