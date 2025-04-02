package test_utils

import "testing"

func Fail(t *testing.T, testNo int, format string, args ...any) {
	t.Errorf("test %d: "+format, append([]any{testNo}, args...)...)
}
