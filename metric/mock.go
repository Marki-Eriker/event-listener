package metric

import "time"

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m Mock) AddPanic(_ bool, _ string) {}

func (m Mock) AddResponse(_ string, _ int, _ time.Duration) {}

func (m Mock) AddDBQuery(_ bool, _ time.Duration, _ string) {}

func (m Mock) AddEvent(_ bool, _ string, _ time.Duration) {}
