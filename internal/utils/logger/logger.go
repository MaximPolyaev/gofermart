package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New(io io.Writer) *Logger {
	lg := logrus.New()
	lg.SetOutput(io)
	lg.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	lg.SetLevel(logrus.InfoLevel)

	return &Logger{Logger: lg}
}
