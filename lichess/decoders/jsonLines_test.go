package decoders

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestStruct struct {
	ID   string  `json:"id"`
	Val1 int64   `json:"val-1"`
	Val2 bool    `json:"val-2"`
	Val3 float64 `json:"val-3"`
}

func TestDecoder_DecodeSimpleNDJson(t *testing.T) {
	const jsonStream = `
{"id": "qwerty", "val-1": 20, "val-2": true, "val-3": 20.00}
{"id": "qwerty1", "val-1": 30, "val-2": false, "val-3": 30.00}
`

	var tests []*TestStruct

	err := NewDecoder(strings.NewReader(jsonStream)).Decode(&tests)

	if err != nil {
		t.Errorf("Decoder returned an error: %v", err)
	}

	want := []*TestStruct{{
		ID:   "qwerty",
		Val1: 20,
		Val2: true,
		Val3: 20.00,
	}, {
		ID:   "qwerty1",
		Val1: 30,
		Val2: false,
		Val3: 30.00,
	}}

	if diff := cmp.Diff(tests, want); diff != "" {
		t.Errorf("Results do not match. Diff: %+v", diff)
	}
}

func TestDecoder_DecodeSingleJson(t *testing.T) {
	const jsonStream = `{"id": "qwerty", "val-1": 20, "val-2": true, "val-3": 20.00}`

	var test *TestStruct

	err := NewDecoder(strings.NewReader(jsonStream)).Decode(&test)

	if err != nil {
		t.Errorf("Decoder returned an error: %v", err)
	}

	want := &TestStruct{
		ID:   "qwerty",
		Val1: 20,
		Val2: true,
		Val3: 20.00,
	}

	if diff := cmp.Diff(test, want); diff != "" {
		t.Errorf("Responses do not match. Diff: %+v", diff)
	}
}

func TestDecoder_ExpectError(t *testing.T) {
	var test *TestStruct

	err := NewDecoder(strings.NewReader("")).Decode(test)

	if err == nil {
		t.Errorf("Account.GetMyEmail returned error: %v", err)
	}
}
