package request

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

const TraceIDHeader = "request-trace-id"

// Request - структура для работы с запросом
type Request struct {
	w         http.ResponseWriter
	r         *http.Request
	beginTime time.Time
	body      []byte
	route     string
	traceID   string
}

// New - создает новый запрос
func New(w http.ResponseWriter, r *http.Request) *Request {
	if responseMetric == nil {
		responseMetric = dummyCallbackResponse
	}

	request := &Request{
		w:         w,
		r:         r,
		beginTime: time.Now(),
	}

	var err error

	if r.Body != nil {
		defer r.Body.Close()
		if request.body, err = io.ReadAll(r.Body); err != nil {
			request.Log().Errorf("Unable to read request body: %s", err)
		} else {
			r.Body = io.NopCloser(bytes.NewReader(request.body))
		}
	}

	muxRoute := mux.CurrentRoute(r)
	if muxRoute != nil {
		request.route, _ = muxRoute.GetPathTemplate()
	}

	request.traceID = r.Header.Get(TraceIDHeader)

	if request.traceID == "" {
		request.traceID = uuid.NewV4().String()
	}

	ctx := context.WithValue(r.Context(), TraceID, request.traceID)
	r = r.WithContext(ctx)

	return request
}

// Log - возвращает обогащенный logger для запроса
func (r *Request) Log() *logrus.Entry {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	return logger.
		WithField("method", r.r.Method).
		WithField("host", r.r.Host).
		WithField("proto", r.r.Proto).
		WithField("remote_addr", r.r.RemoteAddr).
		WithField("request_uri", r.r.RequestURI).
		WithField("route", r.route).
		WithField("duration", time.Now().Sub(r.beginTime))
}

// Query - возвращает query-параметры
func (r *Request) Query() url.Values {
	return r.r.URL.Query()
}

// QueryValue - возвращает по имени аргумент запроса
func (r *Request) QueryValue(name string) Value {
	return Value(r.Query().Get(name))
}

// VarsValue - возвращает по имени переменную пути
func (r *Request) VarsValue(name string) Value {
	return Value(r.GetVar(name))
}

// GetVar - возвращает переменную пути по имени
func (r *Request) GetVar(name string) string {
	return mux.Vars(r.r)[name]
}

// GetBody - возвращает тело запроса в сыром виде
func (r *Request) GetBody() []byte {
	return r.body
}
