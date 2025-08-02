package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/ashprao/ansort"
)

func main() {
	fmt.Println("=== sort.Interface Integration Examples ===")

	// Example 1: Working with Existing Infrastructure
	fmt.Println("\n--- Example 1: Existing Infrastructure Integration ---")
	files := []string{"doc10.txt", "doc2.txt", "doc1.txt", "doc20.txt"}
	fmt.Printf("Files to sort: %v\n", files)

	sorter := ansort.NewSorter(files)
	duration := TimedSort(sorter)
	fmt.Printf("Sorted files: %v\n", files)
	fmt.Printf("Sort duration: %v\n", duration)

	// Example 2: Using Different Sort Algorithms
	fmt.Println("\n--- Example 2: Different Sort Algorithms ---")
	versions := []string{"v1.10.0", "v1.2.0", "v1.20.0", "v1.1.0"}
	fmt.Printf("Versions before: %v\n", versions)

	versionSorter := ansort.NewSorter(versions)

	// Check if already sorted
	if sort.IsSorted(versionSorter) {
		fmt.Println("Versions are already sorted!")
	} else {
		fmt.Println("Versions need sorting...")
		// Use stable sort to maintain order of equal elements
		sort.Stable(versionSorter)
		fmt.Printf("Versions after stable sort: %v\n", versions)
	}

	// Example 3: Performance-Critical Reusable Sorting
	fmt.Println("\n--- Example 3: Performance-Critical Reusable Sorting ---")
	data := []string{"item100", "item10", "item2", "item1", "item50"}
	fmt.Printf("Initial data: %v\n", data)

	// Create a reusable sorter
	reusableSorter := ansort.NewSorter(data, ansort.WithCaseInsensitive())

	fmt.Println("Performing multiple sort operations...")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		// Simulate data modification between sorts
		if i%100 == 0 {
			shuffleSlice(data) // Shuffle every 100 iterations
		}
		sort.Sort(reusableSorter) // Reuse the same sorter object
	}
	totalDuration := time.Since(start)

	fmt.Printf("Final sorted data: %v\n", data)
	fmt.Printf("1000 sorts completed in: %v\n", totalDuration)

	// Example 4: Generic Sorting Framework
	fmt.Println("\n--- Example 4: Generic Sorting Framework ---")
	framework := NewSortingFramework()

	// Add different types of sorters to the framework
	numbers := []string{"file100", "file10", "file2"}
	files2 := []string{"Report10.pdf", "report2.pdf", "REPORT1.pdf"}

	framework.AddSorter("numbers", ansort.NewSorter(numbers))
	framework.AddSorter("files", ansort.NewSorter(files2, ansort.WithCaseInsensitive()))

	fmt.Printf("Before framework sort - numbers: %v\n", numbers)
	fmt.Printf("Before framework sort - files: %v\n", files2)

	framework.SortAll()

	fmt.Printf("After framework sort - numbers: %v\n", numbers)
	fmt.Printf("After framework sort - files: %v\n", files2)

	// Example 5: Benchmarking Different Approaches
	fmt.Println("\n--- Example 5: Benchmarking Comparison ---")
	testData := []string{"file100", "file10", "file2", "file1", "file50", "file5"}
	fmt.Printf("Test data: %v\n", testData)

	// Benchmark convenience function
	testData1 := make([]string, len(testData))
	copy(testData1, testData)
	convenienceDuration := BenchmarkConvenienceFunction(testData1)

	// Benchmark sorter approach
	testData2 := make([]string, len(testData))
	copy(testData2, testData)
	sorterDuration := BenchmarkSorterApproach(testData2)

	fmt.Printf("Convenience function duration: %v\n", convenienceDuration)
	fmt.Printf("Sorter approach duration: %v\n", sorterDuration)
	fmt.Printf("Both results: %v (should be identical)\n", testData1)
}

// TimedSort demonstrates integrating with existing infrastructure
func TimedSort(data sort.Interface) time.Duration {
	start := time.Now()
	sort.Sort(data)
	return time.Since(start)
}

// shuffleSlice simulates data modification between sorts
func shuffleSlice(slice []string) {
	// Simple shuffle - swap first and last elements
	if len(slice) > 1 {
		slice[0], slice[len(slice)-1] = slice[len(slice)-1], slice[0]
	}
}

// SortingFramework demonstrates a generic framework that works with sort.Interface
type SortingFramework struct {
	sorters map[string]sort.Interface
}

func NewSortingFramework() *SortingFramework {
	return &SortingFramework{
		sorters: make(map[string]sort.Interface),
	}
}

func (sf *SortingFramework) AddSorter(name string, sorter sort.Interface) {
	sf.sorters[name] = sorter
}

func (sf *SortingFramework) SortAll() {
	for name, sorter := range sf.sorters {
		fmt.Printf("  Sorting %s (%d items)...\n", name, sorter.Len())
		sort.Sort(sorter)
	}
}

// BenchmarkConvenienceFunction times the convenience function approach
func BenchmarkConvenienceFunction(data []string) time.Duration {
	start := time.Now()
	ansort.SortStrings(data)
	return time.Since(start)
}

// BenchmarkSorterApproach times the sorter + sort.Sort approach
func BenchmarkSorterApproach(data []string) time.Duration {
	start := time.Now()
	sorter := ansort.NewSorter(data)
	sort.Sort(sorter)
	return time.Since(start)
}
