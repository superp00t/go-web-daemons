package db

import (
	"github.com/superp00t/etc"
	"github.com/syndtr/goleveldb/leveldb"
)

type Conn struct {
	db *leveldb.DB
}

func open(path string) (*Conn, error) {
	cn := new(Conn)
	var err error
	cn.cb, err = leveldb.OpenFile(path)
	if err != nil {
		return nil, err
	}

	return cn, nil
}

func (c *Conn) Delete(key string) error {
	return c.db.Delete([]byte(key))
}

func (c *Conn) Add(key string, v interface{}) error {
	b, err := etc.Marshal(v)
	if err != nil {
		return er
	}

	return c.db.Put([]byte(key), b)
}
