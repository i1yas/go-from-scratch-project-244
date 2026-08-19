package diff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeDiff(t *testing.T) {
	userShort := map[string]any{"id": 10, "name": "john"}
	userExtended := map[string]any{"id": "10", "email": "john@john.com", "name": "john"}

	cases := []struct {
		name   string
		value1 map[string]any
		value2 map[string]any
		want   []Item
	}{
		{
			name:   "flat maps",
			value1: map[string]any{"a": -1, "x": 10, "y": 20},
			value2: map[string]any{"a": -1, "x": 100, "z": 300},
			want: []Item{
				{Key: "a", Change: ItemChangeNone, OldValue: -1},
				{Key: "x", Change: ItemChangeReplace, OldValue: 10, NewValue: 100},
				{Key: "y", Change: ItemChangeRemove, OldValue: 20},
				{Key: "z", Change: ItemChangeAdd, NewValue: 300},
			},
		},
		{
			name:   "nested map changed",
			value1: map[string]any{"user": userShort},
			value2: map[string]any{"user": userExtended},
			want: []Item{
				{
					Key:    "user",
					Change: ItemChangeNone,
					Nested: &[]Item{
						{Key: "email", Change: ItemChangeAdd, NewValue: "john@john.com"},
						{Key: "id", Change: ItemChangeReplace, OldValue: 10, NewValue: "10"},
						{Key: "name", Change: ItemChangeNone, OldValue: "john"},
					},
				},
			},
		},
		{
			name:   "nested map unchanged",
			value1: map[string]any{"user": userShort},
			value2: map[string]any{"user": userShort},
			want: []Item{
				{
					Key:    "user",
					Change: ItemChangeNone,
					Nested: &[]Item{
						{Key: "id", Change: ItemChangeNone, OldValue: 10},
						{Key: "name", Change: ItemChangeNone, OldValue: "john"},
					},
				},
			},
		},
		{
			name:   "nested array changed",
			value1: map[string]any{"xs": []any{3, 5, 1}},
			value2: map[string]any{"xs": []any{1, 3, 5, 7}},
			want: []Item{
				{
					Key:      "xs",
					Change:   ItemChangeReplace,
					OldValue: []any{3, 5, 1},
					NewValue: []any{1, 3, 5, 7},
				},
			},
		},
		{
			name:   "nested array unchanged",
			value1: map[string]any{"xs": []any{3, 5, 1}},
			value2: map[string]any{"xs": []any{3, 5, 1}},
			want: []Item{
				{
					Key:      "xs",
					Change:   ItemChangeNone,
					OldValue: []any{3, 5, 1},
				},
			},
		},
		{
			name:   "no change",
			value1: map[string]any{"x": 10},
			value2: map[string]any{"x": 10},
			want: []Item{
				{Key: "x", Change: ItemChangeNone, OldValue: 10},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ComputeDiff(tc.value1, tc.value2)

			require.Equal(t, tc.want, got)
		})
	}
}
