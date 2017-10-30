package gosocketio

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/orisano/golang-socketio/protocol"
)

var (
	ErrorSendTimeout     = errors.New("timeout")
	ErrorSocketOverflood = errors.New("socket overflood")
)

/**
Send message packet to socket
*/
func send(msg *protocol.Message, c *Channel, args interface{}) error {
	if args != nil {
		jargs, err := json.Marshal(args)
		if err != nil {
			return err
		}
		msg.Args = string(jargs)
	}

	command, err := protocol.Encode(msg)
	if err != nil {
		return err
	}
	if len(c.out) == queueBufferSize {
		return ErrorSocketOverflood
	}
	c.out <- command

	return nil
}

/**
Create packet based on given data and send it
*/
func (c *Channel) Emit(method string, args interface{}) error {
	msg := &protocol.Message{
		Type:   protocol.MessageTypeEmit,
		Method: method,
	}

	return send(msg, c, args)
}

/**
Create ack packet based on given data and send it and receive response
*/
func (c *Channel) Ack(method string, args interface{}, timeout time.Duration) (string, error) {
	msg := &protocol.Message{
		Type:   protocol.MessageTypeAckRequest,
		AckId:  c.ack.getNextID(),
		Method: method,
	}

	waiter := make(chan string)
	c.ack.addWaiter(msg.AckId, waiter)

	err := send(msg, c, args)
	if err != nil {
		c.ack.removeWaiter(msg.AckId)
	}

	select {
	case result := <-waiter:
		return result, nil
	case <-time.After(timeout):
		c.ack.removeWaiter(msg.AckId)
		return "", ErrorSendTimeout
	}
}
