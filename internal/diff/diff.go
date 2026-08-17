package diff

import (
	"reflect"
	"slices"
)

// ItemChange represents one of possible changes: add, remove, no change.
type ItemChange int

// ItemChange option
const (
	ItemChangeNone ItemChange = iota
	ItemChangeAdd
	ItemChangeRemove
	ItemChangeReplace
)

// Item is item of internal diff representation
type Item struct {
	Key      string
	OldValue any
	NewValue any
	Change   ItemChange
	Nested   *[]Item
}

// ComputeDiff computes diff and returns internal representation
func ComputeDiff(a, b map[string]any) []Item {
	return computeMapsDiff(a, b)
}

func computeMapsDiff(a, b map[string]any) []Item {
	result := make([]Item, 0, len(a)+len(b))

	keys := getJoinedKeys(a, b)

	slices.Sort(keys)

	for _, k := range keys {
		v1, ok1 := a[k]
		v2, ok2 := b[k]

		if !ok1 {
			result = append(result, Item{
				Key:      k,
				NewValue: v2,
				Change:   ItemChangeAdd,
			})

			continue
		}

		if !ok2 {
			result = append(result, Item{
				Key:      k,
				OldValue: v1,
				Change:   ItemChangeRemove,
			})

			continue
		}

		kind1 := reflect.ValueOf(v1).Kind()
		kind2 := reflect.ValueOf(v2).Kind()

		if kind1 != kind2 {
			result = append(result, Item{
				Key:      k,
				Change:   ItemChangeReplace,
				OldValue: v1,
				NewValue: v2,
			})

			continue
		}

		if kind1 == reflect.Map {
			map1 := v1.(map[string]any)
			map2 := v2.(map[string]any)
			subDiff := computeMapsDiff(map1, map2)
			result = append(result, Item{
				Key:    k,
				Change: ItemChangeNone,
				Nested: &subDiff,
			})

			continue
		}

		if kind1 == reflect.Slice {
			result = append(result, Item{
				Key:      k,
				Change:   ItemChangeReplace,
				OldValue: v1,
				NewValue: v2,
			})

			continue
		}

		if v1 != v2 {
			result = append(result, Item{
				Key:      k,
				OldValue: v1,
				NewValue: v2,
				Change:   ItemChangeReplace,
			})

			continue
		}

		result = append(result, Item{
			Key:      k,
			OldValue: v1,
			Change:   ItemChangeNone,
		})
	}

	return result
}

func getJoinedKeys(a, b map[string]any) []string {
	keys := make([]string, 0, len(a)+len(b))

	for k := range a {
		keys = append(keys, k)
	}

	for k := range b {
		if !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}

	return keys
}
