// +build wasm

package db

import (
	"encoding/json"
	"errors"
	"syscall/js"
)

type Conn struct {
	openReq js.Value
	db      js.Value
	errc    chan error
}

func (c *Conn) onError(args []js.Value) {
	c.errc <- fromError(args[0])
}

func (c *Conn) onSuccess(args []js.Value) {
	c.db = c.openReq.Get("result")
	c.errc <- nil
}

func open(path string) (*Conn, error) {
	cn := new(Conn)
	cn.errc = make(chan error)

	onError := js.NewCallback(cn.onError)
	onSuccess := js.NewCallback(cn.onSuccess)

	cn.openReq = js.Global.Get("indexedDB").Call("open", path)
	cn.openReq.Set("onsuccess", onSuccess)
	cn.openReq.Set("onerror", onError)

	err := <-cn.errc

	if err == nil {
		return cn, nil
	}

	return nil, err
}

func (c *Conn) Add(key string, v interface{}) error {
	table := structName(v)
	tx := c.tx(table)

	objStore := tx.Call("objectStore", table)
	add := objStore.Call("add", key, jsObject(v))

	errc := make(chan error)

	onSuccess := js.NewCallback(func(args []js.Value) {
		errc <- nil
	})

	onError := js.NewCallback(func(args []js.Value) {
		errc <- fromError(args[0])
	})

	add.Set("onsuccess", onSuccess)
	add.Set("onerror", onError)

	err := <-errc

	onSuccess.Release()
	onError.Release()

	return err
}

func (c *Conn) Delete(key string) error {
	table := structName(v)
	tx := c.tx(table)

	objStore := tx.Call("objectStore", table)
	add := objStore.Call("delete", key)

	errc := make(chan error)

	onSuccess := js.NewCallback(func(args []js.Value) {
		errc <- nil
	})

	onError := js.NewCallback(func(args []js.Value) {
		errc <- fromError(args[0])
	})

	add.Set("onsuccess", onSuccess)
	add.Set("onerror", onError)

	err := <-errc

	onSuccess.Release()
	onError.Release()

	return err
}

func (c *Conn) tx(table string) js.Value {
	array := js.Global().Get("Array").New(table)
	tx := js.Call("transaction", array, "readwrite")
	return tx
}

func jsObject(v interface{}) js.Value {
	bytes, _ := json.Marshal(v)
	return js.Global().Get("JSON").Get("parse", string(bytes))
}

func fromError(v js.Value) error {
	return errors.New("db: " + v.String())
}
