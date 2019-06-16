package contract

/*
Imports are exposed to all wasm functions
*/

// #include <stdlib.h>
//
// extern int32_t c_read(void *context);
// extern void c_write(void *context, int32_t ptr);
import "C"

import (
	"fmt"
	"unsafe"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)


//export c_read
func c_read(context unsafe.Pointer) int32 {
	data := ReadDB()
	fmt.Printf("read: %s\n", data)
	return WasmString(data)
}

//export c_write
func c_write(context unsafe.Pointer, ptr int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	text := readString(memory[ptr:])
	fmt.Printf("writing: %s\n", text)
	WriteDB(text)
}

func wasmImports() (*wasm.Imports, error) {
	imp, err := wasm.NewImports().Append("c_read", c_read, C.c_read)
	if err != nil {
		return nil, err
	}
	imp, err = imp.Append("c_write", c_write, C.c_write)
	if err != nil {
		return nil, err
	}
	return imp, nil
}
