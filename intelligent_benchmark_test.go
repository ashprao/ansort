package ansort

import (
	"fmt"
	"sort"
	"testing"
)

// BenchmarkIntelligentSelection verifies that auto-selection provides optimal performance
func BenchmarkIntelligentSelection(t *testing.B) {
	// Generate datasets of different sizes
	generateDataset := func(size int) []string {
		data := make([]string, size)
		for i := 0; i < size; i++ {
			data[i] = fmt.Sprintf("file%d.txt", i*13%1000)
		}
		return data
	}

	// Test small datasets (where legacy should be faster)
	t.Run("Small_Dataset_30", func(b *testing.B) {
		dataset := generateDataset(30)

		b.Run("Intelligent_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStrings(data) // Uses intelligent selection
			}
		})

		b.Run("Always_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStringsLegacy(data)
			}
		})

		b.Run("Always_Cached", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				// Force caching by creating cached sorter directly
				sorter := NewCachedSorter(data)
				sort.Sort(sorter)
			}
		})
	})

	// Test medium datasets (where intelligence should pick the right approach)
	t.Run("Medium_Dataset_150", func(b *testing.B) {
		dataset := generateDataset(150)

		b.Run("Intelligent_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStrings(data) // Uses intelligent selection
			}
		})

		b.Run("Always_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStringsLegacy(data)
			}
		})

		b.Run("Always_Cached", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				sorter := NewCachedSorter(data)
				sort.Sort(sorter)
			}
		})
	})

	// Test large datasets (where caching should always win)
	t.Run("Large_Dataset_500", func(b *testing.B) {
		dataset := generateDataset(500)

		b.Run("Intelligent_Default", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStrings(data) // Uses intelligent selection
			}
		})

		b.Run("Always_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				SortStringsLegacy(data)
			}
		})

		b.Run("Always_Cached", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				data := make([]string, len(dataset))
				copy(data, dataset)
				sorter := NewCachedSorter(data)
				sort.Sort(sorter)
			}
		})
	})

	// Test comparison performance with different string lengths
	t.Run("Compare_Performance", func(b *testing.B) {
		shortStrings := []string{"a1", "a2", "b1", "b2", "c1", "c2"}
		longStrings := []string{
			"very_long_filename_with_many_characters_001.txt",
			"very_long_filename_with_many_characters_002.txt",
			"very_long_filename_with_many_characters_003.txt",
			"very_long_filename_with_many_characters_004.txt",
			"very_long_filename_with_many_characters_005.txt",
			"very_long_filename_with_many_characters_006.txt",
		}

		b.Run("Short_Strings_Intelligent", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(shortStrings); j++ {
					for k := j + 1; k < len(shortStrings); k++ {
						Compare(shortStrings[j], shortStrings[k]) // Uses intelligent selection
					}
				}
			}
		})

		b.Run("Short_Strings_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(shortStrings); j++ {
					for k := j + 1; k < len(shortStrings); k++ {
						CompareLegacy(shortStrings[j], shortStrings[k])
					}
				}
			}
		})

		b.Run("Long_Strings_Intelligent", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(longStrings); j++ {
					for k := j + 1; k < len(longStrings); k++ {
						Compare(longStrings[j], longStrings[k]) // Uses intelligent selection
					}
				}
			}
		})

		b.Run("Long_Strings_Legacy", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(longStrings); j++ {
					for k := j + 1; k < len(longStrings); k++ {
						CompareLegacy(longStrings[j], longStrings[k])
					}
				}
			}
		})
	})
}
