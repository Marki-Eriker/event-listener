package logger

import "github.com/sirupsen/logrus"

// Options - настойки логгера
type Options struct {
	Level       string
	ForceColors bool
}

// New - новый логгер
func New(opt *Options) (*logrus.Logger, error) {
	logLevel, err := logrus.ParseLevel(opt.Level)
	if err != nil {
		return nil, err
	}

	ll := logrus.New()

	ll.SetFormatter(&logrus.TextFormatter{
		ForceColors: opt.ForceColors,
	})

	ll.SetLevel(logLevel)

	return ll, nil
}
