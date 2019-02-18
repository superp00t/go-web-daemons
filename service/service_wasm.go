// +build wasm

package service

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type port struct {
	parent    *Port
	svc       *Service
	port      js.Value
	messageCb js.Callback
	closeCb   js.Callback
}

func (p *port) sendString(s string) {
	p.port.Call("postMessage", s)
}

func (p *port) onMessage(j []js.Value) {
	str := j[0].Get("data").String()
	fmt.Println("got some data", str)

	var i IPC
	err := json.Unmarshal([]byte(str), &i)
	if err == nil {
		p.svc.dispatchIPC(p.parent, i)
	}
}

func (p *port) onClose(j []js.Value) {
	p.messageCb.Release()
	p.closeCb.Release()
}

func (s *Service) onConnect(j []js.Value) {
	p := new(port)
	p.port = j[0].Get("ports").Index(0)

	p.messageCb = js.NewCallback(p.onMessage)
	p.closeCb = js.NewCallback(p.onClose)
	p.port.Call("addEventListener", "message", p.messageCb)
	p.port.Call("addEventListener", "close", p.closeCb)
	p.port.Call("start")

	p.svc = s
	p.parent = &Port{p}
	s.loadPort(p.parent)

	p.parent.PostMessage(IPC{
		Type: "open",
	})
}

func (s *Service) Run() {
	onConnect := js.NewCallback(s.onConnect)
	js.Global().Set("onconnect", onConnect)

	pending := js.Global().Get("pendingConnections")

	for i := 0; i < pending.Length(); i++ {
		pendingConn := pending.Index(i)
		fmt.Println("dispatching queued connection", i)
		s.onConnect([]js.Value{pendingConn})
	}
	<-s.close
}

func (p *Port) Close() {
}
