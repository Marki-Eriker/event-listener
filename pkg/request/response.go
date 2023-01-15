package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FinishOK - завершает запрос удачно с кодом 200
func (r *Request) FinishOK(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusOK).
		WithField("duration", time.Now().Sub(r.beginTime)).
		Infof("Response: %s", fmt.Sprintf(msg, args...))

	r.finish(http.StatusOK, msg, args...)
}

// FinishCreated - завершает запрос удачно с кодом 201
func (r *Request) FinishCreated(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusCreated).
		WithField("duration", time.Now().Sub(r.beginTime)).
		Infof("Response: %s", fmt.Sprintf(msg, args...))

	r.finish(http.StatusOK, msg, args...)
}

// FinishOKJSON - завершает запрос с кодом 200 и объектом для JSON
func (r *Request) FinishOKJSON(i interface{}) {
	r.FinishJSON(http.StatusOK, i)
}

// FinishJSON - завершает запрос с произвольным кодом и объектом для JSON
func (r *Request) FinishJSON(code int, i interface{}) {
	data, err := json.Marshal(i)
	if err != nil {
		r.Log().
			WithField("duration", time.Now().Sub(r.beginTime)).
			Errorf("Unable to marshal response data: %s", err)

		r.FinishError("Unable to marshal response data: %s", err)
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(code)

	if _, err := r.w.Write(data); err != nil {
		r.Log().Warnf("Unable to write data: %s", err)
		return
	}

	ll := r.Log().
		WithField("status", code).
		WithField("duration", time.Now().Sub(r.beginTime))

	if code < 300 {
		ll.Info("Response")
	} else if code >= 300 && code < 500 {
		ll.Warn("Response")
	} else {
		ll.Error("Response")
	}

	responseMetric(r.route, code, time.Since(r.beginTime))
}

// FinishError - завершает запрос неудачно с кодом 500
func (r *Request) FinishError(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusInternalServerError).
		WithField("duration", time.Now().Sub(r.beginTime)).
		Errorf("Response: %s", fmt.Sprintf(msg, args...))

	r.finish(http.StatusInternalServerError, msg, args...)
}

// FinishBadRequest - завершает запрос неудачно с кодом 400
func (r *Request) FinishBadRequest(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusBadRequest).
		WithField("duration", time.Now().Sub(r.beginTime)).
		Warnf("Response: %s", fmt.Sprintf(msg, args...))

	r.finish(http.StatusBadRequest, msg, args...)
}

// FinishUnauthorized - завершает запрос неудачно с кодом 401
func (r *Request) FinishUnauthorized(msg string, args ...interface{}) {
	r.Log().
		WithField("status", http.StatusUnauthorized).
		WithField("duration", time.Now().Sub(r.beginTime)).
		Warnf("Response: %s", fmt.Sprintf(msg, args...))

	r.finish(http.StatusUnauthorized, msg, args...)
}

func (r *Request) finish(code int, msg string, args ...interface{}) {
	r.w.WriteHeader(code)
	buf := bytes.NewBufferString(fmt.Sprintf(msg, args...))
	r.w.Write(buf.Bytes())
	responseMetric(r.route, code, time.Since(r.beginTime))
}
