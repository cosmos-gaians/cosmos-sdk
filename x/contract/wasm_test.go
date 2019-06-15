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

	res, err := run(simple, "add1", []interface{}{int32(7), int32(9)}, AsString)
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

	initMsg := `{
		"contract_address": "deadbeef",
		"sender": "0123456789",
		"sent_funds": 1000,
		"msg": {
			"verifier": "ethan",
			"beneficiary": "jehan"
		}
	}`

	res, err := Run(regen, "init", []interface{}{initMsg})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(res.Msgs) != 0 {
		t.Fatalf("Unexpected result: %v", res)
	}

	badSend := `{
		invalid: 123
	}`

	res, err = Run(regen, "send", []interface{}{badSend})
	if err == nil {
		t.Fatal("Allowed bad json")
	}

	unauthSend := `{
		"contract_address": "deadbeef",
		"sender": "0123456789",
		"sent_funds": 50,
		"msg": {}
	}`

	res, err = Run(regen, "send", []interface{}{unauthSend})
	if err == nil {
		t.Fatal("Allowed no auth")
	}

	goodSend := `{
		"contract_address": "deadbeef",
		"sender": "ethan",
		"sent_funds": 50,
		"msg": {}
	}`

	res, err = Run(regen, "send", []interface{}{goodSend})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	// This is placeholder on success
	if len(res.Msgs) != 1 {
		t.Fatalf("Unexpected result: %v", res)
	}
}
