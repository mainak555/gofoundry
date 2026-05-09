package loggers

import (
	"sync"

	"github.com/go-kit/log"
)

type KitLogWriters struct {
	mu      sync.Mutex
	loggers []log.Logger
}

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

func NewKitLogWriters(loggers ...log.Logger) log.Logger {
	return &KitLogWriters{loggers: loggers}
}
