package gosocketio

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrorWaiterNotFound = errors.New("waiter not found")
)

/**
Processes functions that require answers, also known as acknowledge or ack
*/
type ackProcessor struct {
	counter       int64
	resultWaiters *sync.Map
}

func newAckProcessor() *ackProcessor {
	return &ackProcessor{
		counter:       0,
		resultWaiters: new(sync.Map),
	}
}

/**
get next id of ack call
*/
func (a *ackProcessor) getNextID() int64 {
	return atomic.AddInt64(&a.counter, 1)
}

/**
Just before the ack function called, the waiter should be added
to wait and receive response to ack call
*/
func (a *ackProcessor) addWaiter(id int64, w chan string) {
	a.resultWaiters.Store(id, w)
}

/**
removes waiter that is unnecessary anymore
*/
func (a *ackProcessor) removeWaiter(id int64) {
	a.resultWaiters.Delete(id)
}

/**
check if waiter with given ack id is exists, and returns it
*/
func (a *ackProcessor) getWaiter(id int64) (chan string, error) {
	if waiter, ok := a.resultWaiters.Load(id); ok {
		return waiter.(chan string), nil
	}
	return nil, ErrorWaiterNotFound
}
