package processor

import (
	"github.com/marki-eriker/event-listener/db"
	"github.com/marki-eriker/event-listener/pkg/encoder"
	"github.com/marki-eriker/event-listener/pkg/hasher"
	"github.com/marki-eriker/event-listener/pkg/token"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	log     *logrus.Logger
	store   *Store
	token   token.Generator
	hasher  hasher.Hasher
	encoder encoder.Encoder
}

func New(
	log *logrus.Logger,
	store *Store,
	token token.Generator,
	hasher hasher.Hasher,
	encoder encoder.Encoder,
) *Processor {
	return &Processor{
		log:     log,
		store:   store,
		token:   token,
		hasher:  hasher,
		encoder: encoder,
	}
}

type Store struct {
	User  db.User
	Event db.Event
}

func (p *Processor) Log(traceID string) *logrus.Entry {
	return p.log.WithField("trace_id", traceID)
}
