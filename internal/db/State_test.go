//go:guild integration

package db_test

import (
	db "expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/internal/test_utils"
	"expat-news/queue-manager/pkg/utils"
	"testing"
)

func compare_test_db_states(a, b *db.State) bool {
	if a == nil && a != b {
		return false
	}
	if *a.Id != *b.Id {
		return false
	}
	if *a.Name != *b.Name {
		return false
	}
	return true
}

func TestStateGet(t *testing.T) {
	test_utils.SetupDB(t)
	tests := []struct {
		id       *int
		name     *string
		expected db.State
	}{
		{
			utils.Ptr(1),
			nil,
			db.State{utils.Ptr(1), utils.Ptr("new")},
		},
		{
			nil,
			utils.Ptr("active"),
			db.State{utils.Ptr(2), utils.Ptr("active")},
		},
		{
			utils.Ptr(0),
			utils.Ptr("done"),
			db.State{utils.Ptr(0), utils.Ptr("done")},
		},
	}
	for i, test := range tests {
		st := &db.State{test.id, test.name}
		if err := st.Get(); err != nil {
			t.Errorf("test %d: unexpected error occured: %v", i, err)
			continue
		}
		if !compare_test_db_states(&test.expected, st) {
			t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, st)
			continue
		}
	}
}
