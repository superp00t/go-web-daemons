package service

import (
	"encoding/json"
	"sync"
)

type IPC struct {
	Type  string          `json:"type"`
	Query string          `json:"query,omitempty"`
	ID    string          `json:"id,omitempty"`
	Data  json.RawMessage `json:"data"`
}

type Service struct {
	onPort func(p *Port)

	ports  []*Port
	portsL *sync.Mutex

	query *sync.Map
	close chan bool
}

type Port struct {
	*port
}

func (s *Service) OnPort(fn func(p *Port)) {
	s.onPort = fn
}

func New() *Service {
	s := new(Service)
	s.portsL = new(sync.Mutex)
	s.close = make(chan bool)
	s.query = new(sync.Map)
	return s
}

func (s *Service) loadPort(p *Port) {
	s.portsL.Lock()
	s.ports = append(s.ports, p)
	s.portsL.Unlock()

	if s.onPort != nil {
		go s.onPort(p)
	}
}

func (p *Port) PostMessage(obj IPC) {
	str, _ := json.Marshal(obj)
	p.sendString(string(str))
}

func (svc *Service) Broadcast(jso interface{}) {
	iface, err := json.Marshal(jso)
	if err != nil {
		panic(err)
	}

	obj := IPC{
		Type: "broadcast",
		Data: json.RawMessage(iface),
	}

	str, _ := json.Marshal(obj)

	s := string(str)
	svc.portsL.Lock()
	for _, v := range svc.ports {
		v.sendString(s)
	}
	svc.portsL.Unlock()
}

type Query struct {
	Port *Port
	ID   string
	Data json.RawMessage
}

func (q *Query) Send(v interface{}) {
	dat, _ := json.Marshal(v)
	q.Port.PostMessage(IPC{
		Type: "answer",
		ID:   q.ID,
		Data: dat,
	})
}

func (svc *Service) On(query string, fn func(q *Query)) {
	svc.query.Store(query, fn)
}

func (svc *Service) dispatchIPC(port *Port, i IPC) {
	switch i.Type {
	case "query":
		fn, ok := svc.query.Load(i.Query)
		if !ok {
			return
		}

		f := fn.(func(q *Query))
		q := &Query{
			Port: port,
			ID:   i.ID,
			Data: i.Data,
		}

		go f(q)
	case "decommission":
		svc.portsL.Lock()
		idx := -1
		for i, v := range svc.ports {
			if v == port {
				idx = i
				break
			}
		}

		if idx == 1 {
			panic("tried to delete port that was not in slice")
		}

		svc.ports = append(svc.ports[:idx], svc.ports[idx+1:]...)
		svc.portsL.Unlock()

		port.Close()
	}
}
