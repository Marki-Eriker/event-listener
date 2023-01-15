package api

import (
	"github.com/marki-eriker/event-listener/metric"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Listener struct {
	AddEventMetric func(success bool, from string, duration time.Duration)
	Addr           string
	log            *logrus.Logger
	lis            net.Listener
	handlerEvent   func(net.Conn)
	stop           chan struct{}
}

type EventListenerOption struct {
	Addr string
}

func NewEventListener(
	opt *EventListenerOption,
	log *logrus.Logger,
	m metric.Metric,
	handler func(net.Conn),
) *Listener {
	return &Listener{
		AddEventMetric: m.AddEvent,
		Addr:           opt.Addr,
		log:            log,
		handlerEvent:   handler,
	}
}

func (l *Listener) Start() error {
	lis, err := net.Listen("tcp", l.Addr)
	if err != nil {
		return err
	}

	l.lis = lis

	l.log.Warnf("starting listener on %s", l.Addr)

	for {
		conn, err := lis.Accept()
		if err != nil {
			l.log.Warn("connection close")
			break
		}

		l.log.WithField("from", conn.RemoteAddr().String()).Info("New TCP connection")

		go l.handlerEvent(conn)
	}

	return nil
}

func (l *Listener) Stop() error {
	return l.lis.Close()
}
