package wq

func Move[T any](arr []T, src, dest int) []T {
	if src < 0 || dest < 0 {
		// If src or dest are negative, return array as is.
		return arr
	}

	if src >= len(arr) || dest >= len(arr) {
		// If src or dest are out of bounds, return array as is.
		return arr
	}

	if src == dest {
		// If src and dest are the same, no need to move anything.
		return arr
	}

	result := make([]T, len(arr))
	for i := range arr {
		// do not change anything on the left of the region
		if i < src && i < dest {
			result[i] = arr[i]
			continue
		}

		// do not change anything on the right of the region
		if i > src && i > dest {
			result[i] = arr[i]
			continue
		}

		// move source item to destination
		if i == src {
			result[dest] = arr[i]
			continue
		}

		if src < dest {
			result[i-1] = arr[i]
			continue
		}

		if src > dest {
			result[i+1] = arr[i]
			continue
		}
	}

	return result
}
