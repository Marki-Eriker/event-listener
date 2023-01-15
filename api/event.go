package api

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/marki-eriker/event-listener/entity/event"
	"github.com/marki-eriker/event-listener/entity/paginator"
	"github.com/marki-eriker/event-listener/pkg/request"
	"github.com/marki-eriker/event-listener/processor"
	"io"
	"net"
	"net/http"
	"time"
)

type eventsOutput struct {
	Paginator *paginator.Output `json:"paginator"`
	Events    []*event.Event    `json:"events"`
}

func (s *Server) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	req := request.New(w, r)

	page := req.QueryValue("page").MustUInt()
	pageSize := req.QueryValue("page_size").MustUInt()
	eventID := req.QueryValue("event_id").MustUInt()
	SystemName := req.QueryValue("system_name").String()
	OrderField := req.QueryValue("order_field").String()
	Order := req.QueryValue("order").String()

	pr := paginator.NewInput(page, pageSize, eventID, SystemName, OrderField, Order)

	events, total, err := s.processor.GetEvents(r.Context(), pr)

	switch err {
	case nil:
		res := &eventsOutput{
			Paginator: paginator.NewOutput(pr, total),
			Events:    events,
		}

		req.FinishOKJSON(res)
	case request.ErrRoleNotFoundInContext:
		req.FinishUnauthorized("")
	case processor.ErrUnexpectedDatabaseBehavior:
		req.FinishError(err.Error())
	default:
		req.FinishError("unexpected error")
	}

}

func (s *Server) handleMarkEventAnIncident(w http.ResponseWriter, r *http.Request) {
	req := request.New(w, r)

	id := req.VarsValue("id").MustUUID()
	if id == nil {
		req.FinishBadRequest("no valid id (uuid) in request")
		return
	}

	err := s.processor.MarkEventAnIncident(r.Context(), id)

	switch err {
	case nil:
		req.FinishOK("event was successfully marked")
	case processor.ErrUnexpectedDatabaseBehavior:
		req.FinishError(err.Error())
	default:
		req.FinishError("unexpected error")
	}

}

func (s *Server) HandleEvent(conn net.Conn) {
	r := bufio.NewReader(conn)
	from := conn.RemoteAddr().String()

	for {
		bt := time.Now()
		b, rErr := r.ReadBytes('\n')
		if rErr != nil && rErr != io.EOF {
			break
		}

		var e event.IncomingPayload

		err := json.Unmarshal(b, &e)
		if err != nil {
			s.log.WithError(err).Error("Unable to unmarshal JSON data")

			if rErr == io.EOF {
				break
			}

			continue
		}

		err = s.processor.ProcessingEvent(context.Background(), &e)
		s.metric.AddEvent(err == nil, from, time.Since(bt))

		if rErr == io.EOF {
			break
		}
	}
}
