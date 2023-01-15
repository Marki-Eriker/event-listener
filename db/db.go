package db

import (
	"context"
	c "github.com/kaatinga/const-errs"
	"github.com/marki-eriker/event-listener/entity/event"
	"github.com/marki-eriker/event-listener/entity/paginator"
	"github.com/marki-eriker/event-listener/entity/user"
	uuid "github.com/satori/go.uuid"
)

const (
	ErrRecordNotFound c.Error = "record not found"
	ErrDuplicateKey   c.Error = "entity with same key already exists"
)

type Event interface {
	Insert(context.Context, *event.Event) error
	GetMany(ctx context.Context, p *paginator.Input, withSecured bool) (e []*event.Event, total uint, err error)
	UpdateIncident(context.Context, *uuid.UUID, bool) error
}

type User interface {
	GetByID(context.Context, uint) (*user.User, error)
	GetByLogin(context.Context, string) (*user.User, error)
	Insert(context.Context, *user.User) error
	UpdateVerify(context.Context, uint, bool) error
}
