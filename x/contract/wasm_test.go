package contract

import (
	"strings"
	"testing"
)

func TestImportFunc(t *testing.T) {
	simple, err := Read("examples/import_func/build/import_func.wasm")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	res, err := Run(simple, "add1", []interface{}{int32(7), int32(9)}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if res.(string) != strings.Repeat("fool ", 17) {
		t.Fatalf("Unexpected result: %d", res)
	}
}
