package testutils

import (
	"bytes"
	"encoding/json"
	"testing"
)

func AssertStructEq(t *testing.T, exp, got interface{}) {
	bexp, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		t.Fatalf("failed to enocode payload: %v", exp)
	}

	bgot, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("failed to enocode payload: %v", got)
	}

	res := bytes.Compare(bexp, bgot)
	if res != 0 {
		t.Errorf("exp = %s", string(bexp[:]))
		t.Errorf("got = %s", string(bgot[:]))
		t.Fatal("exp and got structs are not equivalent")
	}
}
