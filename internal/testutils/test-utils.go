package testutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	MONGODB_URI_TEST      string = "mongodb://localhost:27017"
	MONGODB_DATABASE_TEST string = "ctb-unit-tests"
)

func AssertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic but none occurred")
		}
	}()
	f()
}

func AssertTrue(t *testing.T, value bool, comp string) {
	t.Helper()
	if !value {
		t.Errorf("%s | exp=true, got=false", comp)
	}
}

func AssertFalse(t *testing.T, value bool, comp string) {
	t.Helper()
	if value {
		t.Errorf("%s | exp=false, got=true", comp)
	}
}

func AssertNil(t *testing.T, value interface{}, comp string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s | exp=nil, got=%v", comp, format(value))
		return
	}
}

func AssertNotNil(t *testing.T, value interface{}, comp string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s | exp!=nil, got=nil", comp)
	}
}

func AssertEq(t *testing.T, exp, got interface{}, comp string) {
	t.Helper()
	if exp == nil && got == nil {
		return
	}

	if exp != nil && got == nil {
		t.Errorf("%s | exp=%v, got=nil", comp, format(exp))
		return
	}

	if exp == nil && got != nil {
		t.Errorf("%s | exp=nil, got=%v", comp, format(got))
		return
	}

	kexp := reflect.ValueOf(exp).Kind()
	kgot := reflect.ValueOf(got).Kind()
	if kexp != kgot {
		t.Errorf("%s | exp_type=%v, got_type=%v", comp, kexp, kgot)
	}

	exp = format(exp)
	got = format(got)
	if fmt.Sprintf("%v", exp) == fmt.Sprintf("%v", got) {
		return
	}

	t.Errorf("%s | exp=%v, got=%v", comp, exp, got)
}

func format(value interface{}) interface{} {
	if value == nil {
		return "nil"
	}

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		payload := reflect.Indirect(reflect.ValueOf(value))
		reflect.ValueOf(value).Elem().Set(payload)
	}

	if reflect.ValueOf(value).Kind() == reflect.Struct {
		bytes, _ := json.MarshalIndent(value, "", "  ")
		return string(bytes[:])
	}

	if reflect.ValueOf(value).Kind() == reflect.Map {
		bytes, _ := json.MarshalIndent(value, "", "  ")
		return string(bytes[:])
	}

	if reflect.ValueOf(value).Kind() == reflect.Slice {
		bytes, _ := json.MarshalIndent(value, "", "  ")
		return string(bytes[:])
	}

	return value
}
