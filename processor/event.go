package processor

import (
	"context"
	"fmt"
	"github.com/marki-eriker/event-listener/entity/event"
	"github.com/marki-eriker/event-listener/entity/paginator"
	"github.com/marki-eriker/event-listener/entity/user"
	"github.com/marki-eriker/event-listener/pkg/recovery"
	"github.com/marki-eriker/event-listener/pkg/request"
	uuid "github.com/satori/go.uuid"
)

// ProcessingEvent - обработает входящие событие и сохранит в БД
func (p *Processor) ProcessingEvent(ctx context.Context, payload *event.IncomingPayload) error {
	defer recovery.Recover("processor.ProcessingEvent")

	ll := p.Log(request.GetTraceID(ctx)).
		WithField("event_id", payload.EventID).
		WithField("system", payload.SystemName)

	ll.Debug("ProcessingEvent begin")
	defer ll.Debug("ProcessingEvent end")

	e := event.FromPayload(payload, p.encoder)

	err := p.store.Event.Insert(ctx, e)
	if err != nil {
		ll.Errorf("unable to save event: %s", err)
		return err
	}

	return nil
}

// GetEvents - вернет события по заданным параметрам фильтрации.
// Может вернуть request.ErrRoleNotFoundInContext, ErrUnexpectedDatabaseBehavior
func (p *Processor) GetEvents(ctx context.Context, pr *paginator.Input) (e []*event.Event, total uint, err error) {
	defer recovery.Recover("processor.GetEvents")

	ll := p.Log(request.GetTraceID(ctx)).WithFields(pr.LogFields())

	ll.Debug("GetEvents begin")
	defer ll.Debug("GetEvents end")

	role, err := request.GetUserRole(ctx)
	if err != nil {
		return nil, 0, err
	}

	fmt.Println(role)

	e, total, err = p.store.Event.GetMany(ctx, pr, role == user.Analyst)
	if err != nil {
		ll.Errorf("unable to get events: %s", err)
		return nil, 0, ErrUnexpectedDatabaseBehavior
	}

	return e, total, nil
}

// MarkEventAnIncident - отметит событие как инцидент.
// Может вернуть ErrUnexpectedDatabaseBehavior
func (p *Processor) MarkEventAnIncident(ctx context.Context, id *uuid.UUID) error {
	defer recovery.Recover("processor.MarkEventAnIncident")

	ll := p.Log(request.GetTraceID(ctx)).WithField("id", id)

	ll.Debug("MarkEventAnIncident begin")
	defer ll.Debug("MarkEventAnIncident end")

	err := p.store.Event.UpdateIncident(ctx, id, true)
	if err != nil {
		ll.Errorf("unable to update event: %s", err)
		return ErrUnexpectedDatabaseBehavior
	}

	return nil
}
