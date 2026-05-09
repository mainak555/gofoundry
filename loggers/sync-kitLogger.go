package loggers

import (
	"sync"

	"github.com/go-kit/log"
)

// KitLogWriters fan-outs log events to multiple synchronized log.Logger instances.
type KitLogWriters struct {
	mu      sync.Mutex
	loggers []log.Logger
}

// Log writes a single log event to every configured logger.
func (sw *KitLogWriters) Log(keyvals ...interface{}) error {
	defer sw.mu.Unlock()
	sw.mu.Lock()
	for _, w := range sw.loggers {
		err := w.Log(keyvals...)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewKitLogWriters builds a thread-safe logger that forwards to loggers.
func NewKitLogWriters(loggers ...log.Logger) log.Logger {
	return &KitLogWriters{loggers: loggers}
}
