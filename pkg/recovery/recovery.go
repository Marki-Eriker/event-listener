package recovery

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

// PanicMetrics - метрика паники
type PanicMetrics func(restart bool, module string)

var (
	panicMetricsHandler PanicMetrics
)

func Setup(metricsHandler PanicMetrics) {
	panicMetricsHandler = metricsHandler

	if panicMetricsHandler == nil {
		panicMetricsHandler = dummyPanicMetrics
	}
}

// Recover - восстанавливает приложение
func Recover(module string) {
	err := recover()
	if err != nil {
		panicMetricsHandler(false, module)
		logrus.WithError(fmt.Errorf("%#v", err)).Errorf("Panic: %s", debug.Stack())
	}
}

// Restart - перезапускает приложение
func Restart(module string) {
	err := recover()
	if err != nil {
		panicMetricsHandler(true, module)
		logrus.WithError(fmt.Errorf("%#v", err)).Fatalf("Panic (restart): %s", debug.Stack())
	}
}

func dummyPanicMetrics(_ bool, _ string) {}
