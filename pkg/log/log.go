package log

import (
	"github.com/go-logr/logr"
	"sync"
	"time"
)

var (
	Log  logr.Logger = dLog
	dLog             = NewDelegatingLogger(logr.Discard())

	loggerWasSetLock sync.Mutex
	loggerWasSet     bool
)

func SetLogger(logger logr.Logger) {
	loggerWasSetLock.Lock()
	defer loggerWasSetLock.Unlock()

	loggerWasSet = true
	dLog.Fulfill(logger)
}

func init() {
	go func() {
		time.Sleep(30 * time.Second)
		loggerWasSetLock.Lock()
		defer loggerWasSetLock.Unlock()
		if !loggerWasSet {
			dLog.Fulfill(logr.Discard())
		}
	}()
}
