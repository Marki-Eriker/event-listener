package postgres

import (
	"context"
	_ "embed"
	"github.com/go-pg/pg/v10"
	"time"
)

//go:embed migrate.sql
var migrateQuery string

type Postgres struct {
	User    *UserStore
	Event   *EventStore
	storage *pg.DB
}

type Options struct {
	Addr     string
	User     string
	Password string
	Database string
}

func New(opt *Options, metric func(bool, time.Duration, string)) (*Postgres, error) {
	storage := pg.Connect(&pg.Options{
		Addr:     opt.Addr,
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
	})

	err := storage.Ping(context.Background())

	if err != nil {
		return nil, err
	}

	err = migrate(storage)
	if err != nil {
		return nil, err
	}

	us := &UserStore{storage: storage, queryMetric: metric}
	es := &EventStore{storage: storage, queryMetric: metric}

	return &Postgres{User: us, Event: es, storage: storage}, err
}

func (pdb *Postgres) Stop() error {
	return pdb.storage.Close()
}

func migrate(storage *pg.DB) error {
	_, err := storage.Exec(migrateQuery)
	return err
}
