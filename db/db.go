package db

import "reflect"

func Open(path string) *Conn {
	return open(path)
}

func structName(v interface{}) string {
	return reflect.TypeOf(v).Name()
}
