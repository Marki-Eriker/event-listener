package request

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	logger         *logrus.Logger
	responseMetric func(string, int, time.Duration)
	initializer    sync.Once
)

// Setup - функция устанавливает логгер и метрику ответов
func Setup(l *logrus.Logger, respMetric func(string, int, time.Duration)) {
	initializer.Do(func() {
		logger = l
		responseMetric = respMetric
	})
}

func dummyCallbackResponse(_ string, _ int, _ time.Duration) {}
