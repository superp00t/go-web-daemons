// +build wasm

package ws

import (
	"fmt"
	"io"
	"syscall/js"
	"time"

	"github.com/superp00t/go-web-daemons/wasmutils"
)

func prototype() js.Value {
	return js.Global().Get("WebSocket")
}

type conn struct {
	binary       bool
	conn         js.Value
	dialComplete chan bool
	recv         chan []byte
	errc         chan error
	u8Proto      js.Value
}

func dialOpts(o Opts) (*conn, error) {
	if o.Socks5 != "" {
		return nil, fmt.Errorf("ws: cannot dial WebSocket in WebAssembly using a SOCKS5 proxy parameter")
	}

	o.binary = o.Binary

	args := []interface{}{o.URL}
	if o.Subprotocol != "" {
		args = append(args, o.Subprotocol)
	}

	c := &conn{}
	c.dialComplete = make(chan bool)
	c.errc = make(chan error)
	c.recv = make(chan []byte)
	c.conn = prototype().New(args...)
	c.u8Proto = js.Global().Get("Uint8Array")

	c.conn.Call("addEventListener", "open", js.NewCallback(func(a []js.Value) {
		c.dialComplete <- true
	}))

	c.conn.Call("addEventListener", "close", js.NewCallback(func(a []js.Value) {
		c.errc <- io.EOF
	}))

	c.conn.Call("addEventListener", "error", js.NewCallback(func(a []js.Value) {
		c.errc <- fmt.Errorf("ws: ", a[0].String())
	}))

	bl := js.Global().Get("Blob")
	ab := js.Global().Get("ArrayBuffer")

	c.conn.Call("addEventListener", "message", js.NewCallback(func(a []js.Value) {
		event := a[0]
		data := event.Get("data")

		if data.InstanceOf(ab) {
			c.handleArrayBuffer(data)
		} else if data.InstanceOf(bl) {
			c.handleBlob(data)
		} else {
			c.recv <- []byte(data.String())
		}
	}))

	<-c.dialComplete
	return c, nil
}

func (c *conn) Send(b []byte) error {
	if c.binary {
		u8 := js.TypedArrayOf(b)
		c.conn.Call("send", u8)
		u8.Release()
	} else {
		str := js.ValueOf(string(b))
		c.conn.Call("send", str)
	}

	return nil
}

func (c *conn) Recv() ([]byte, error) {
	select {
	case err := <-c.errc:
		return nil, err
	case b := <-c.recv:
		return b, nil
	}
}

func (c *conn) Close() error {
	c.conn.Call("close")
	return nil
}

func (c *conn) handleArrayBuffer(ab js.Value) {
	fmt.Println("arraybuffer")
	out := wasmutils.LoadBytesFromArrayBuffer(ab)
	fmt.Println(time.Since(t))
	c.recv <- out
}

func (c *conn) handleBlob(blob js.Value) {
	fr := js.Global().Get("FileReader").New()
	fr.Set("onload", js.NewCallback(func(a []js.Value) {
		c.handleArrayBuffer(a[0].Get("target").Get("result"))
	}))
	fr.Call("readAsArrayBuffer", blob)
}
