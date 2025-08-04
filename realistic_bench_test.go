package ansort

import (
	"fmt"
	"testing"
)

// BenchmarkRealisticPerformance tests with larger, more realistic datasets
func BenchmarkRealisticPerformance(t *testing.B) {
	// Generate a more realistic dataset - larger and with repeated patterns
	generateLargeDataset := func(size int) []string {
		data := make([]string, size)
		patterns := []string{
			"file%d.txt", "document%d.pdf", "image%d.jpg", "video%d.mp4",
			"backup%d.zip", "config%d.yaml", "log%d.txt", "data%d.csv",
		}

		for i := 0; i < size; i++ {
			pattern := patterns[i%len(patterns)]
			// Mix different number ranges to create realistic sorting scenarios
			num := (i * 13) % 1000 // Creates numbers like 0, 13, 26, 39, ..., 987, 0, 13...
			data[i] = fmt.Sprintf(pattern, num)
		}
		return data
	}

	// Test with small dataset (cache overhead may dominate)
	t.Run("Small_Dataset_100", func(b *testing.B) {
		dataset := generateLargeDataset(100)

		b.Run("SortStrings_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStrings(data)
			}
		})

		b.Run("SortStrings_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStringsLegacy(data)
			}
		})
	})

	// Test with medium dataset (cache should start helping)
	t.Run("Medium_Dataset_1000", func(b *testing.B) {
		dataset := generateLargeDataset(1000)

		b.Run("SortStrings_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStrings(data)
			}
		})

		b.Run("SortStrings_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStringsLegacy(data)
			}
		})
	})

	// Test repeated sorting (where cache really shines)
	t.Run("Repeated_Sorting_500x10", func(b *testing.B) {
		dataset := generateLargeDataset(500)

		b.Run("SortStrings_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Sort the same dataset 10 times (simulating repeated operations)
				for j := 0; j < 10; j++ {
					data := make([]string, len(dataset))
					copy(data, dataset)
					SortStrings(data)
				}
			}
		})

		b.Run("SortStrings_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Sort the same dataset 10 times
				for j := 0; j < 10; j++ {
					data := make([]string, len(dataset))
					copy(data, dataset)
					SortStringsLegacy(data)
				}
			}
		})
	})

	// Test pure comparison performance (where caching helps most)
	t.Run("Pure_Comparisons", func(b *testing.B) {
		dataset := generateLargeDataset(100)

		b.Run("Compare_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Do many comparisons between same strings (cache helps here)
				for j := 0; j < len(dataset); j++ {
					for k := j + 1; k < len(dataset); k++ {
						Compare(dataset[j], dataset[k])
					}
				}
			}
		})

		b.Run("Compare_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Same comparisons with legacy implementation
				for j := 0; j < len(dataset); j++ {
					for k := j + 1; k < len(dataset); k++ {
						CompareLegacy(dataset[j], dataset[k])
					}
				}
			}
		})
	})
}
