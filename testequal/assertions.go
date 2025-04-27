//go:build !solution

package testequal

import (
	"fmt"
	"reflect"
)

func makeErr(t T, msgAndArgs []interface{}) {
	t.Helper()
	if len(msgAndArgs) == 0 {
		t.Errorf("")
		t.FailNow()
	}
	err := msgAndArgs[0].(string)
	t.Errorf(fmt.Sprintf(err, msgAndArgs[1:]...))
	t.FailNow()
}

func makeErr2(t T, msgAndArgs []interface{}) bool {
	t.Helper()
	if len(msgAndArgs) == 0 {
		return false
	}
	err := msgAndArgs[0].(string)
	t.Errorf(fmt.Sprintf(err, msgAndArgs[1:]...))
	return false
}

func isEqualSlices[S int | string | byte](slice1, slice2 []S) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	if len(slice1) == 0 && len(slice2) == 0 {
		if slice1 == nil && slice2 != nil || slice2 == nil && slice1 != nil {
			return false
		}
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func isEqualMap(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	if len(map1) == 0 && len(map2) == 0 {
		if map1 == nil && map2 != nil || map2 == nil && map1 != nil {
			return false
		}
	}
	for key, val1 := range map1 {
		val2, ok := map2[key]
		if !ok || val1 != val2 {
			return false
		}
	}

	return true
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)
	if expectedType != actualType {
		return false
	}
	switch expected.(type) {
	case int, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		if actual != expected {
			return makeErr2(t, msgAndArgs)
		}
	case []int:
		exp := expected.([]int)
		act := actual.([]int)
		if !isEqualSlices(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	case []byte:
		exp := expected.([]byte)
		act := actual.([]byte)
		if !isEqualSlices(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	case map[string]string:
		exp := expected.(map[string]string)
		act := actual.(map[string]string)
		if !isEqualMap(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	default:
		return false
	}

	return true
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)
	if expectedType != actualType {
		return true
	}
	switch expected.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		if actual == expected {
			return makeErr2(t, msgAndArgs)
		}
	case []int:
		exp := expected.([]int)
		act := actual.([]int)
		if isEqualSlices(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	case []byte:
		exp := expected.([]byte)
		act := actual.([]byte)
		if isEqualSlices(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	case map[string]string:
		exp := expected.(map[string]string)
		act := actual.(map[string]string)
		if isEqualMap(exp, act) {
			return makeErr2(t, msgAndArgs)
		}
	}
	return true
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)
	if expectedType != actualType {
		t.Errorf("")
		t.FailNow()
	}
	switch expected.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		if actual != expected {
			makeErr(t, msgAndArgs)
		}
	case []int:
		exp := expected.([]int)
		act := actual.([]int)
		if !isEqualSlices(exp, act) {
			makeErr(t, msgAndArgs)
		}
	case []byte:
		exp := expected.([]byte)
		act := actual.([]byte)
		if !isEqualSlices(exp, act) {
			makeErr(t, msgAndArgs)
		}
	case map[string]string:
		exp := expected.(map[string]string)
		act := actual.(map[string]string)
		if !isEqualMap(exp, act) {
			makeErr(t, msgAndArgs)
		}
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)
	if expectedType == actualType {
		switch expected.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
			if actual == expected {
				makeErr(t, msgAndArgs)
			}
		case []int:
			exp := expected.([]int)
			act := actual.([]int)
			if isEqualSlices(exp, act) {
				makeErr(t, msgAndArgs)
			}
		case []byte:
			exp := expected.([]byte)
			act := actual.([]byte)
			if isEqualSlices(exp, act) {
				makeErr(t, msgAndArgs)
			}
		case map[string]string:
			exp := expected.(map[string]string)
			act := actual.(map[string]string)
			if isEqualMap(exp, act) {
				makeErr(t, msgAndArgs)
			}
		}
	}
}
