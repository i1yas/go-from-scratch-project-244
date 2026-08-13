package diff

import (
	"reflect"
	"slices"
)

// ItemChange represents one of possible changes: add, remove, no change.
type ItemChange int

// ItemChange option
const (
	ItemChangeNone = iota
	ItemChangeAdd
	ItemChangeRemove
)

// Item is item of internal diff representation
type Item struct {
	key    string
	value  any
	change ItemChange
	nested *[]Item
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
				key:    k,
				value:  v2,
				change: ItemChangeAdd,
			})

			continue
		}

		if !ok2 {
			result = append(result, Item{
				key:    k,
				value:  v1,
				change: ItemChangeRemove,
			})

			continue
		}

		kind1 := reflect.ValueOf(v1).Kind()
		kind2 := reflect.ValueOf(v2).Kind()

		if kind1 != kind2 {
			result = append(result, Item{
				key:    k,
				value:  v1,
				change: ItemChangeRemove,
			})

			result = append(result, Item{
				key:    k,
				value:  v2,
				change: ItemChangeAdd,
			})

			continue
		}

		if kind1 == reflect.Map {
			map1 := v1.(map[string]any)
			map2 := v2.(map[string]any)
			subDiff := computeMapsDiff(map1, map2)
			result = append(result, Item{
				key:    k,
				nested: &subDiff,
			})

			continue
		}

		if kind1 == reflect.Slice {
			slice1 := v1.([]any)
			slice2 := v2.([]any)
			subDiff := computeSlicesDiff(slice1, slice2)
			result = append(result, Item{
				key:    k,
				nested: &subDiff,
			})

			continue
		}

		if v1 != v2 {
			result = append(result, Item{
				key:    k,
				value:  v1,
				change: ItemChangeRemove,
			})

			result = append(result, Item{
				key:    k,
				value:  v2,
				change: ItemChangeAdd,
			})

			continue
		}

		result = append(result, Item{
			key:    k,
			value:  v1,
			change: ItemChangeNone,
		})
	}

	return result
}

func computeSlicesDiff(a, b []any) []Item {
	result := make([]Item, 0, len(a)+len(b))

	// TODO

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
