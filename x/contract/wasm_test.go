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

func TestRegenInit(t *testing.T) {
	regen, err := Read("examples/regen/build/regen.wasm")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// Set up global static for test
	data = "{}"

	json := `{
		"sender": "0123456789",
		"init_funds": 1000,
		"init_msg": {
			"verifier": "ethan",
			"beneficiary": "jehan"
		}
	}`

	res, err := Run(regen, "init", []interface{}{json}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if res.(string) != "" {
		t.Fatalf("Unexpected result: %d", res)
	}
}
