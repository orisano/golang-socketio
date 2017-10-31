package gosocketio

import (
	"net/url"

	"github.com/orisano/golang-socketio/transport"
)

/**
Socket.io client representation
*/
type Client struct {
	Methods *methods
	Channel *Channel
}

/**
Get ws/wss url by host and port
*/
func GetURL(host string, secure bool, spath ...string) string {
	u := new(url.URL)
	if secure {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	p := "/socket.io/"
	if len(spath) == 1 {
		p = spath[0]
	}

	u.Host = host
	u.Path = p
	q := u.Query()
	q.Set("EIO", "3")
	q.Set("transport", "websocket")
	u.RawQuery = q.Encode()

	return u.String()
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, tr transport.Transport) (*Client, error) {
	conn, err := tr.Connect(url)
	if err != nil {
		return nil, err
	}
	c := &Client{
		Channel: NewChannel(conn),
		Methods: newMethods(),
	}

	go inLoop(c.Channel, c.Methods)
	go outLoop(c.Channel, c.Methods)
	go pinger(c.Channel)

	return c, nil
}

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(c.Channel, c.Methods)
}
