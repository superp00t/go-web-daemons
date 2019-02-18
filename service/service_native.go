// +build !wasm

package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/superp00t/etc"

	"github.com/gorilla/websocket"
)

type port struct {
	conn   *websocket.Conn
	closed bool
}

func (p *port) sendString(data string) {
	p.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (s *Service) Run() {
	serviceToken := etc.GenerateRandomUUID()
	mx := http.NewServeMux()

	mx.HandleFunc("/ipc", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != serviceToken.String() {
			return
		}

		ws := websocket.Upgrader{}
		c, err := ws.Upgrade(rw, r, nil)
		if err != nil {
			return
		}

		port := &Port{&port{c, false}}
		if s.onPort != nil {
			go s.onPort(port)
		}

		for {
			var i IPC
			_, msg, err := c.ReadMessage()
			if err != nil {
				port.Close()
				return
			}

			err = json.Unmarshal(msg, &i)
			if err != nil {
				port.Close()
				return
			}

			s.dispatchIPC(port, i)
		}
	})

	var err error
	var l net.Listener

	for x := 2000; x < 65535; x++ {
		host := fmt.Sprintf("localhost:%d", x)
		uri := "ws://" + host + "/ipc?token=" + serviceToken.String()

		l, err = net.Listen("tcp", host)
		if err != nil {
			continue
		}

		fmt.Println(uri)

		err = http.Serve(l, mx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	fmt.Println(err)
	os.Exit(1)
}

func (p *port) Close() {
	if !p.closed {
		p.conn.Close()
		p.closed = true
	}
}
