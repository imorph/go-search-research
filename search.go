package search

func LinearSearchFloat64s(haystack []float64, needle float64) int {
	n := len(haystack)
	i := 0

	// Unroll the loop for better instruction-level parallelism
	for n >= 4 {
		if haystack[i] >= needle {
			if haystack[i] == needle {
				return i
			}
			return len(haystack)
		}
		if haystack[i+1] >= needle {
			if haystack[i+1] == needle {
				return i + 1
			}
			return len(haystack)
		}
		if haystack[i+2] >= needle {
			if haystack[i+2] == needle {
				return i + 2
			}
			return len(haystack)
		}
		if haystack[i+3] >= needle {
			if haystack[i+3] == needle {
				return i + 3
			}
			return len(haystack)
		}
		i += 4
		n -= 4
	}

	// Handle the remaining elements
	for n > 0 {
		if haystack[i] >= needle {
			if haystack[i] == needle {
				return i
			}
			return len(haystack)
		}
		i++
		n--
	}
	return len(haystack)
}
