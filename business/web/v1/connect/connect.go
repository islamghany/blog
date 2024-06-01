package connect

import (
	"context"
	"github/islamghany/blog/foundation/logger"
	"time"
)

// ConnectWithBackOff will try to connect to the service with backoff mechanism
// it will try to connect to the service with the given function
func ConnectWithBackOff[T any](ctx context.Context, log *logger.Logger, connectName string, fn func() (*T, error), backOff int) (*T, error) {
	timeToSleep := 1
	var connection *T
	for {
		c, err := fn()
		if err == nil {
			connection = c
			break
		}
		backOff--
		if backOff == 0 {
			log.Error(ctx, connectName, "connection error", err)
			return nil, err
		}
		log.Info(ctx, connectName, "error", err, "retrying in", timeToSleep, "seconds")
		timeToSleep = timeToSleep * backOff * 2
		log.Info(ctx, connectName, "retrying in", timeToSleep, "seconds")
		time.Sleep(time.Duration(timeToSleep) * time.Second)
		continue
	}
	return connection, nil
}
