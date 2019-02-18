// +build !wasm

package wasmutils

func init() {
	panic("you cannot import wasmutils in native executables.")
}
