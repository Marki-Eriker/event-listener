package metric

import "time"

type Metric interface {
	AddPanic(restart bool, source string)
	AddResponse(route string, status int, duration time.Duration)
	AddDBQuery(success bool, duration time.Duration, query string)
	AddEvent(success bool, from string, duration time.Duration)
}
