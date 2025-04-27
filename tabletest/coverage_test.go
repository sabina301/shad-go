//go:build !change

package tabletest

import (
	"errors"
	"testing"
	"time"
)

// min coverage: . 95%

type TestCase struct {
	req    string
	expDur time.Duration
	expErr error
}

func TestTable(t *testing.T) {
	for _, tc := range []TestCase{
		{"0", 0, nil},
		{"-", 0, errors.New("time: invalid duration -")},
	} {
		t.Run("", func(t *testing.T) {
			actDur, actErr := ParseDuration(tc.req)
			if actDur != tc.expDur || !errors.Is(actErr, tc.expErr) {
				t.Errorf("failed to parse duration %q", tc.req)
			}
		})
	}
}
