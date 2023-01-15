package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/marki-eriker/event-listener/metric"
	"github.com/marki-eriker/event-listener/pkg/token"
	"github.com/marki-eriker/event-listener/processor"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/pprof"
	"time"
)

type Server struct {
	router          *mux.Router
	httpsrv         *http.Server
	processor       *processor.Processor
	log             *logrus.Logger
	shutdownTimeout time.Duration
	addr            string
	metric          metric.Metric
	tokenManager    token.Generator
}

type HTTPAPIOptions struct {
	Addr            string
	Debug           bool
	ShutdownTimeout time.Duration
}

func NewHTTPAPI(
	p *processor.Processor,
	log *logrus.Logger,
	opt *HTTPAPIOptions,
	m metric.Metric,
	t token.Generator,
) *Server {
	s := &Server{
		router:       mux.NewRouter(),
		processor:    p,
		log:          log,
		addr:         opt.Addr,
		metric:       m,
		tokenManager: t,
	}

	s.initMainRoutesV1()

	if opt.Debug {
		s.initDebugRoutes()
	}

	lr := s.router
	s.httpsrv = &http.Server{
		Handler: lr,
		Addr:    opt.Addr,
	}

	return s
}

// Serve - запуск HTTP API сервера
func (s *Server) Serve() error {
	s.log.Warnf("starting http api on %s", s.addr)

	return s.httpsrv.ListenAndServe()
}

// Stop - остановка работы HTTP API сервера
func (s *Server) Stop() error {
	if s.shutdownTimeout == 0 {
		return s.httpsrv.Close()
	}

	ctx, cFunc := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cFunc()

	return s.httpsrv.Shutdown(ctx)
}

func (s *Server) initDebugRoutes() {
	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/heap", pprof.Index)
	s.router.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
