package gosocketio

import "testing"

func TestGetURL(t *testing.T) {
	ts := []struct {
		host     string
		secure   bool
		spath    string
		expected string
	}{
		{
			host:     "ws.example",
			secure:   false,
			expected: "ws://ws.example/socket.io/?EIO=3&transport=websocket",
		},
		{
			host:     "ws.example",
			secure:   true,
			expected: "wss://ws.example/socket.io/?EIO=3&transport=websocket",
		},
		{
			host:     "ws.example:8888",
			secure:   false,
			expected: "ws://ws.example:8888/socket.io/?EIO=3&transport=websocket",
		},
		{
			host:     "ws.example",
			secure:   true,
			spath:    "/custompath/",
			expected: "wss://ws.example/custompath/?EIO=3&transport=websocket",
		},
	}

	for _, tc := range ts {
		var got string
		if len(tc.spath) == 0 {
			got = GetURL(tc.host, tc.secure)
		} else {
			got = GetURL(tc.host, tc.secure, tc.spath)
		}
		if got != tc.expected {
			t.Errorf("unexpected url. expected: %v, but got: %v", tc.expected, got)
		}
	}
}
