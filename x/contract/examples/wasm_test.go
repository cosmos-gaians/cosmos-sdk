package examples

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

	initMsg := `{
		"sender": "0123456789",
		"init_funds": 1000,
		"init_msg": {
			"verifier": "ethan",
			"beneficiary": "jehan"
		}
	}`

	res, err := Run(regen, "init", []interface{}{initMsg}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if res.(string) != "" {
		t.Fatalf("Unexpected result: %d", res)
	}

	badSend := `{
		invalid: 123
	}`

	res, err = Run(regen, "send", []interface{}{badSend}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	// Returns "Could not parse MsgSendContract json." on parse error
	if res.(string) != "Could not parse MsgSendContract json." {
		t.Fatalf("Unexpected result: %d", res)
	}

	unauthSend := `{
		"sender": "0123456789",
		"payment": 20,
		"msg": "TODO"
	}`

	res, err = Run(regen, "send", []interface{}{unauthSend}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	// Returns "" on error
	if res.(string) != "" {
		t.Fatalf("Unexpected result: %d", res)
	}

	goodSend := `{
		"sender": "ethan",
		"payment": 20,
		"msg": "TODO"
	}`

	res, err = Run(regen, "send", []interface{}{goodSend}, AsString)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	// This is placeholder on success
	if res.(string) != "Send tx goes here !!!" {
		t.Fatalf("Unexpected result: %d", res)
	}

}
