package search

import (
	"fmt"
	"sort"
	"testing"
)

var result int

func generateSortedFloat64s(n int) []float64 {
	haystack := make([]float64, n)
	for i := range haystack {
		haystack[i] = float64(i)
	}
	return haystack
}

func BenchmarkSearchFunctions(b *testing.B) {
	lengths := []int{10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000}
	positions := []string{"beginning", "middle", "end", "notfound"}

	for _, n := range lengths {
		haystack := generateSortedFloat64s(n)
		for _, pos := range positions {
			var needle float64
			switch pos {
			case "beginning":
				index := n / 10 // Near the beginning
				if index == 0 {
					index = 0
				}
				needle = haystack[index]
			case "middle":
				index := n / 2 // Middle position
				needle = haystack[index]
			case "end":
				index := n - n/10 - 1 // Near the end
				if index < 0 {
					index = n - 1
				}
				needle = haystack[index]
			case "notfound":
				needle = -1 // Value not in haystack
			}

			b.Run(fmt.Sprintf("Linear/n=%d/pos=%s", n, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = LinearSearchFloat64s(haystack, needle)
				}
			})

			b.Run(fmt.Sprintf("Binary/n=%d/pos=%s", n, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = sort.SearchFloat64s(haystack, needle)
				}
			})
		}
	}
}