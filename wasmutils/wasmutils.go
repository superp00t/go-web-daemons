// +build wasm

package wasmutils

import (
	"syscall/js"
)

var uint8Array *js.Value

func acquireU8Constructor() js.Value {
	if uint8Array == nil {
		u8 := js.Global().Get("Uint8Array")
		uint8Array = &u8
	}

	return *uint8Array
}

// Copies bytes from ArrayBuffer value into []byte
func LoadBytesFromArrayBuffer(ab js.Value) []byte {
	u8 := acquireU8Constructor()

	ln := ab.Get("byteLength").Int()
	b := make([]byte, ln)
	write := js.TypedArrayOf(b)

	read := u8.New(ab)
	write.Call("set", read)
	write.Release()

	return b
}

//
// var memory *js.Value
// var membuffer *js.Value
//
// func AcquireMemoryBuffer() js.Value {
// 	u8 := acquireU8Constructor()

// 	if memory == nil {
// 		// Load memory reference
// 		const nanHeader = 0x7FF80000
// 		var memPtr uint64 = nanHeader<<32 | uint64(6)
// 		ref := js.Value{}

// 		rf := reflect.ValueOf(&ref).Elem().Field(0)
// 		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()

// 		rf.SetUint(memPtr)
// 		memory = &ref

// 		write := u8.New(ref.Get("buffer"))
// 		membuffer = &write
// 	}

// 	return *membuffer
// }
//
//// Criminally bad function that is necessary to get good performance.
// func LoadBytesFromArrayBuffer(ab js.Value) []byte {
// write := AcquireMemoryBuffer()

// ln := ab.Get("byteLength").Int()

// // Create memory address to store bytes
// b := make([]byte, ln)
// pointer := uintptr(unsafe.Pointer(&b[0]))
// offset := js.ValueOf(pointer)

// // Copy bytes from array buffer into []byte
// read := u8.New(ab)
// write.Call("set", read, offset)
// return b
//}
