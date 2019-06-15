package contract

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

var (
	curInstance *wasm.Instance
)

// Read loads a wasm file
func Read(filename string) ([]byte, error) {
	return wasm.ReadBytes(filename)
}

// WasmString can be called by a go function provided into Imports
// It will allocate space in wasm, copy the string there, and return a pointer
// The pointer can be returned to the wasm caller to receive the string
func WasmString(res string) int32 {
	return prepareString(*curInstance, res)
}

type ResultParser func(wasm.Instance, wasm.Value) (interface{}, error)

func AsInt32(_ wasm.Instance, res wasm.Value) (interface{}, error) {
	return res.ToI32(), nil
}

func AsInt64(_ wasm.Instance, res wasm.Value) (interface{}, error) {
	return res.ToI64(), nil
}

func AsString(instance wasm.Instance, res wasm.Value) (interface{}, error) {
	outputPointer := res.ToI32()

	memory := instance.Memory.Data()[outputPointer:]
	nth := 0
	var output strings.Builder
	for {
		if memory[nth] == 0 {
			break
		}

		output.WriteByte(memory[nth])
		nth++
	}

	// Deallocate the subject, and the output.
	deallocate, ok := instance.Exports["deallocate"]
	if ok {
		lengthOfOutput := nth
		deallocate(outputPointer, lengthOfOutput)
	}

	// TODO
	// deallocate(inputPointer, lengthOfSubject)

	return output.String(), nil
}

// Run will execute the named function on the wasm bytes with the passed arguments.
// Returns the result or an error
func Run(code []byte, imports *wasm.Imports, call string, args []interface{}, parse ResultParser) (interface{}, error) {
	if imports == nil {
		imports = wasm.NewImports()
	}

	// Instantiates the WebAssembly module.
	instance, err := wasm.NewInstanceWithImports(code, imports)
	if err != nil {
		return nil, errors.Wrap(err, "init wasmer")
	}

	// we give access to some globals for go callbacks
	curInstance = &instance
	defer func() {
		instance.Close()
		curInstance = nil
	}()

	f, ok := instance.Exports[call]
	if !ok {
		return nil, errors.Errorf("Function %s not in Exports", call)
	}

	fArgs := prepareArgs(instance, args)

	ret, err := f(fArgs...)
	if err != nil {
		return nil, errors.Wrap(err, "Execution failure")
	}
	fmt.Printf("%v: %v\n", ret.GetType(), ret)

	return parse(instance, ret)
}

func prepareArgs(instance wasm.Instance, args []interface{}) []interface{} {
	out := make([]interface{}, len(args))

	for i, arg := range args {
		switch t := arg.(type) {
		case int32, int64:
			out[i] = arg
		case string:
			out[i] = prepareString(instance, t)
		case []byte:
			out[i] = prepareString(instance, string(t))
		default:
			panic(fmt.Sprintf("Unsupported type: %T", arg))
		}
	}
	return out
}

func prepareString(instance wasm.Instance, arg string) int32 {
	l := len(arg)
	allocateResult, _ := instance.Exports["allocate"](l)
	inputPointer := allocateResult.ToI32()

	// Write the subject into the memory.
	memory := instance.Memory.Data()[inputPointer:]
	copy(memory, arg)

	// C-string terminates by NULL.
	memory[l] = 0

	return inputPointer
}
