package notification

import "sync/atomic"

type eventHandler func() bool
type eventCallback func(bool)

type subscription struct {
	id             string
	handlerResult  atomic.Value
	resubscribe    bool
	eventHandler   eventHandler
	eventCallbacks []eventCallback
}

func (sub *subscription) setResult(r bool) {
	sub.handlerResult.Store(r)
}

func (sub *subscription) getLatestResult() bool {
	r := sub.handlerResult.Load()
	if r == nil {
		return false
	}
	return r.(bool)
}
