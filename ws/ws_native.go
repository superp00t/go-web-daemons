// +build !wasm

package ws

import (
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

type conn struct {
	conn  *websocket.Conn
	l     *sync.Mutex
	mtype int
}

func (c *conn) Close() error {
	return c.conn.Close()
}

func (c *conn) Send(b []byte) error {
	c.l.Lock()
	err := c.conn.WriteMessage(c.mtype, []byte(b))
	c.l.Unlock()
	return err
}

func (c *conn) Recv() ([]byte, error) {
	_, b, err := c.conn.ReadMessage()
	return b, err
}

func dialOpts(d Opts) (*conn, error) {
	url := d.URL
	c := &conn{}
	var ws *websocket.Conn
	c.l = new(sync.Mutex)

	var dialer = websocket.DefaultDialer

	if d.Socks5 != "" {
		netDialer, err := proxy.SOCKS5("tcp", d.Socks5, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		dialer = &websocket.Dialer{NetDial: netDialer.Dial}
	}

	ws, _, err := dialer.Dial(url, d.Header())
	if err != nil {
		return nil, err
	}

	c.mtype = websocket.TextMessage
	if d.Binary {
		c.mtype = websocket.BinaryMessage
	}

	c.conn = ws

	return c, nil
}
