package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/marki-eriker/event-listener/db"
	"github.com/marki-eriker/event-listener/entity/event"
	"github.com/marki-eriker/event-listener/entity/paginator"
	uuid "github.com/satori/go.uuid"
	"time"
)

type EventStore struct {
	storage     *pg.DB
	queryMetric func(success bool, duration time.Duration, query string)
}

func (es *EventStore) Insert(ctx context.Context, e *event.Event) error {
	bt := time.Now()

	_, err := es.storage.ModelContext(ctx, e).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			es.queryMetric(false, time.Since(bt), "insert_event")
			return db.ErrDuplicateKey
		}

		es.queryMetric(false, time.Since(bt), "insert_event")
		return err
	}

	es.queryMetric(true, time.Since(bt), "insert_event")
	return nil
}

func (es *EventStore) GetMany(ctx context.Context, pr *paginator.Input, withSecured bool) ([]*event.Event, uint, error) {
	bt := time.Now()

	var events []*event.Event

	limit := pr.PageSize
	offset := (pr.Page - 1) * pr.PageSize

	var order string
	if pr.OrderField == "id" || pr.OrderField == "event_id" {
		order = fmt.Sprintf("%s %s", pr.OrderField, pr.Order)
	} else {
		order = fmt.Sprintf("%s %s", "created", pr.Order)
	}

	q := es.storage.ModelContext(ctx, &events)

	if !withSecured {
		q.Column("id", "event_id", "created", "system_name", "incident")
	}

	if pr.SystemName != "" {
		q.Where("system_name = ?", pr.SystemName)
	}

	if pr.EventID != 0 {
		q.Where("event_id = ?", pr.EventID)
	}

	q.Order(order).Limit(int(limit)).Offset(int(offset))

	count, err := q.SelectAndCount()
	if err != nil {
		es.queryMetric(false, time.Since(bt), "get_events")
		return nil, 0, err
	}

	es.queryMetric(true, time.Since(bt), "get_events")
	return events, uint(count), nil
}

func (es *EventStore) UpdateIncident(ctx context.Context, id *uuid.UUID, incident bool) error {
	bt := time.Now()

	e := &event.Event{
		ID:       *id,
		Incident: incident,
	}

	_, err := es.storage.ModelContext(ctx, e).Column("incident").WherePK().Update()
	if err != nil {
		es.queryMetric(false, time.Since(bt), "incident_event")
		return err
	}

	es.queryMetric(true, time.Since(bt), "incident_event")
	return nil
}
