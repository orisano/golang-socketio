package gosocketio

import (
	"sync"

	"github.com/orisano/golang-socketio/protocol"
)

const (
	OnConnection    = "connection"
	OnDisconnection = "disconnection"
	OnError         = "error"
)

/**
System handler function for internal event processing
*/
type Callback func(c *Channel, args string) interface{}

/**
Contains maps of message processing functions
*/
type methods struct {
	messageHandlers     map[string]Callback
	messageHandlersLock sync.RWMutex
}

func newMethods() *methods {
	return &methods{
		messageHandlers: make(map[string]Callback),
	}
}

/**
Add message processing function, and bind it to given method
*/
func (m *methods) On(method string, callback Callback) {
	m.messageHandlersLock.Lock()
	defer m.messageHandlersLock.Unlock()
	m.messageHandlers[method] = callback
}

/**
Find message processing function associated with given method
*/
func (m *methods) findMethod(method string) (Callback, bool) {
	m.messageHandlersLock.RLock()
	defer m.messageHandlersLock.RUnlock()

	f, ok := m.messageHandlers[method]
	return f, ok
}

func (m *methods) callLoopEvent(c *Channel, event string) {
	f, ok := m.findMethod(event)
	if !ok {
		return
	}
	f(c, "")
}

/**
Check incoming message
On ack_resp - look for waiter
On ack_req - look for processing function and send ack_resp
On emit - look for processing function
*/
func (m *methods) processIncomingMessage(c *Channel, msg *protocol.Message) {
	switch msg.Type {
	case protocol.MessageTypeEmit:
		f, ok := m.findMethod(msg.Method)
		if !ok {
			return
		}
		f(c, msg.Args)

	case protocol.MessageTypeAckRequest:
		f, ok := m.findMethod(msg.Method)
		if !ok {
			return
		}
		r := f(c, msg.Args)
		if r == nil {
			r = struct{}{}
		}
		ack := &protocol.Message{
			Type:  protocol.MessageTypeAckResponse,
			AckId: msg.AckId,
		}
		send(ack, c, r)

	case protocol.MessageTypeAckResponse:
		waiter, err := c.ack.getWaiter(msg.AckId)
		if err == nil {
			waiter <- msg.Args
		}
	}
}
