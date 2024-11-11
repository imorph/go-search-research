package search

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sort"
	"sync"
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

// generateRandomFloat64s generates a sorted slice of random float64s of the given length.
func generateRandomSortedFloat64s(length int) []float64 {
	s := make([]float64, length)
	for i := range s {
		s[i] = rand.NormFloat64() + 50.1
	}
	sort.Float64s(s)
	return s
}

var resultFindBucket int

const haystackLen = 90

func BenchmarkSearchFloat64sRandom(b *testing.B) {
	length := haystackLen // Adjust the length of the haystack as needed.
	haystack := generateRandomSortedFloat64s(length)
	mean := (slices.Max(haystack) + slices.Min(haystack)) / 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.NormFloat64() + mean

		// Perform the search.
		resultFindBucket = sort.SearchFloat64s(haystack, needle)
	}
}

func BenchmarkSearchLinearRandom(b *testing.B) {
	length := haystackLen // Adjust the length of the haystack as needed.
	haystack := generateRandomSortedFloat64s(length)
	mean := (slices.Max(haystack) + slices.Min(haystack)) / 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle := rand.NormFloat64() + mean

		// Perform the search.
		resultFindBucket = BasicLinearSearchFloat64s(haystack, needle)
	}
}

var needle float64

func BenchmarkRandomNeedleOverhead(b *testing.B) {
	length := 5000 // Adjust the length of the haystack as needed.
	haystack := generateRandomSortedFloat64s(length)
	mean := (slices.Max(haystack) + slices.Min(haystack)) / 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		needle = rand.NormFloat64() + mean
	}
}

func BenchmarkSearchFunctions(b *testing.B) {
	lengths := []int{10, 20, 30, 35, 40, 50, 60, 100}
	positions := []string{"beginning", "middle", "end", "too_low", "too_high"}

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
			case "too_low":
				needle = -1.0 // Value not in haystack
			case "too_high":
				needle = 500_000_000.0
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

func BenchmarkParallelSearches(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8, 16}
	lengths := []int{10, 20, 30, 35, 40, 50, 60, 100}
	positions := []string{"beginning", "middle", "end", "too_low", "too_high"}

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
			case "too_low":
				needle = -1.0 // Value not in haystack
			case "too_high":
				needle = 500_000_000.0
			}

			for _, c := range concurrencyLevels {
				b.Run(fmt.Sprintf("Linear/n=%d/pos=%s/conc=%d", n, pos, c), func(b *testing.B) {
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						var wg sync.WaitGroup
						wg.Add(c)
						for j := 0; j < c; j++ {
							go func() {
								defer wg.Done()
								result = LinearSearchFloat64s(haystack, needle)
							}()
						}
						wg.Wait()
					}
				})

				b.Run(fmt.Sprintf("Binary/n=%d/pos=%s/conc=%d", n, pos, c), func(b *testing.B) {
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						var wg sync.WaitGroup
						wg.Add(c)
						for j := 0; j < c; j++ {
							go func() {
								defer wg.Done()
								result = sort.SearchFloat64s(haystack, needle)
							}()
						}
						wg.Wait()
					}
				})
			}
		}
	}
}

func BenchmarkLinearSearchImplementations(b *testing.B) {
	lengths := []int{10, 20, 30, 35, 40, 50, 60, 100}
	positions := []string{"beginning", "middle", "end", "too_low", "too_high"}

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
			case "too_low":
				needle = -1.0 // Value not in haystack
			case "too_high":
				needle = 500_000_000.0
			}

			b.Run(fmt.Sprintf("OptimizedLinear/n=%d/pos=%s", n, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = LinearSearchFloat64s(haystack, needle)
				}
			})

			b.Run(fmt.Sprintf("BasicLinear/n=%d/pos=%s", n, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = BasicLinearSearchFloat64s(haystack, needle)
				}
			})

			b.Run(fmt.Sprintf("BinarySearch/n=%d/pos=%s", n, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = sort.SearchFloat64s(haystack, needle)
				}
			})
		}
	}
}
