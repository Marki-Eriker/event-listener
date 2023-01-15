package app

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	signalChannel chan os.Signal
	exitSignal    chan bool
	log           *logrus.Logger
}

// New - новое приложение
func New(log *logrus.Logger) *Application {
	return &Application{log: log}
}

// Run - запускает приложение
func (a *Application) Run() {
	a.log.Warn("Starting application")

	a.exitSignal = make(chan bool)

	a.log.Warn("Application running, waiting for exit signal")

	go a.initSignals()
	<-a.exitSignal
}

func (a *Application) initSignals() {
	a.signalChannel = make(chan os.Signal, 1)
	signal.Notify(a.signalChannel, syscall.SIGTERM)
	signal.Notify(a.signalChannel, syscall.SIGINT)
	signal.Notify(a.signalChannel, syscall.SIGKILL)
	for {
		select {
		case s := <-a.signalChannel:
			switch s {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				close(a.signalChannel)
				a.log.Warnf("We got %s, shutdown application...", s)
				a.exitSignal <- true
				return
			}
		}
	}
}
