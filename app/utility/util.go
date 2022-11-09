package utility

import (
	"context"
	"errors"
	"fmt"
	"keeper/app/pkg/logger"
	"runtime/debug"
	"time"
)

func Retry(name string, f func() bool, times int, interval time.Duration) bool {
	return RetryWithCtx(context.Background(), name, f, times, interval)
}

func SyncRetry(name string, f func() bool, times int, interval time.Duration) {
	for i := 0; i < times; i++ {
		defer func() {
			if r := recover(); r != nil {
				logger.Warnf("Recover from %s retry[%d/%d]: %v. Stacktrace: %s", name, i+1, times, r, debug.Stack())
			}
		}()
		if f() {
			return
		} else {
			logger.Infof("The %d/%d try for %s failed, sleep %d seconds and try again", i+1, times, name, interval/time.Second)
			time.Sleep(interval)
		}
	}
}

func RetryWithCtx(ctx context.Context, name string, f func() bool, times int, interval time.Duration) bool {
	for i := 0; i < times; i++ {
		done := make(chan bool)
		go func(i int) {
			logger.Debugf("Try %s the %d/%d time", name, i+1, times)
			defer func() {
				if r := recover(); r != nil {
					logger.Warnf("Recover from %s retry[%d/%d]: %v. Stacktrace: %s", name, i+1, times, r, debug.Stack())
					done <- false
				}
				close(done)
			}()
			if f() {
				done <- true
			} else {
				logger.Infof("The %d/%d try for %s failed, sleep %d seconds and try again", i+1, times, name, interval/time.Second)
				done <- false
			}
		}(i)
		if <-done {
			logger.Debugf("%s successfully done at the %d/%d try", name, i+1, times)
			return true
		}
		if i == times-1 {
			// Avoid sleep if the last retry has been failed
			logger.Infof("All tries for %s failed ", name)
			return false
		}
		logger.Infof("The %d/%d try for %s failed, sleep %d seconds and try again", i+1, times, name, interval/time.Second)
		select {
		case <-ctx.Done():
			logger.Warnf("The %d/%d try for %s failed due to context done ", i+1, times, name)
			return false
		case <-time.After(interval):
		}
	}

	return false
}

func WithRecover(fn func(), panicHandler func(err error)) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			switch t := r.(type) {
			case string:
				err = errors.New(fmt.Sprintf("err: %s \n %s", t, stack))
			case error:
				err = errors.New(fmt.Sprintf("err: %s \n %s", t.Error(), stack))
			default:
				err = fmt.Errorf("Recover result: %v \n %s", r, stack)
			}
			panicHandler(err)
		}
	}()

	fn()
}
