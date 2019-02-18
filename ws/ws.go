package ws

import "net/http"

type Opts struct {
	URL         string
	Subprotocol string
	Binary      bool
	Socks5      string
	Origin      string
	UserAgent   string
}

type Conn struct {
	*conn
}

func (o Opts) Header() http.Header {
	h := http.Header{}

	if o.Origin != "" {
		h.Set("Origin", o.Origin)
	}

	if o.UserAgent != "" {
		h.Set("User-Agent", o.UserAgent)
	}

	return h
}

func DialOpts(o Opts) (*Conn, error) {
	c, err := dialOpts(o)
	return &Conn{c}, err
}
