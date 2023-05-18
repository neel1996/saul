package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

func NewLogger(ctx ...context.Context) *logrus.Entry {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if len(ctx) == 0 {
		return logrus.New().WithContext(context.Background())
	}

	return logrus.New().WithContext(ctx[0])
}
