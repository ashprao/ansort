package main

import (
	"fmt"
	"sort"

	"github.com/ashprao/ansort"
)

func main() {
	fmt.Println("=== Phase 2.1: Case Sensitivity Demo ===")

	data := []string{"File2.txt", "file10.txt", "FILE1.txt", "item2", "Item10", "ITEM1"}

	// Case-sensitive sorting (default)
	fmt.Println("\n--- Case-Sensitive Sorting (Default) ---")
	caseSensitive := make([]string, len(data))
	copy(caseSensitive, data)
	ansort.SortStrings(caseSensitive)
	fmt.Printf("Before: %v\n", data)
	fmt.Printf("After:  %v\n", caseSensitive)

	// Case-insensitive sorting
	fmt.Println("\n--- Case-Insensitive Sorting ---")
	caseInsensitive := make([]string, len(data))
	copy(caseInsensitive, data)
	ansort.SortStrings(caseInsensitive, ansort.WithCaseInsensitive())
	fmt.Printf("Before: %v\n", data)
	fmt.Printf("After:  %v\n", caseInsensitive)

	// Demonstrating Compare with functional options
	fmt.Println("\n--- String Comparison Examples ---")
	fmt.Printf("Compare(\"File1.txt\", \"file1.txt\") with case-sensitive = %d\n",
		ansort.Compare("File1.txt", "file1.txt")) // default is case-sensitive
	fmt.Printf("Compare(\"File1.txt\", \"file1.txt\") with case-insensitive = %d\n",
		ansort.Compare("File1.txt", "file1.txt", ansort.WithCaseInsensitive()))

	// Demonstrating NewSorter with Go's standard sort.Sort
	fmt.Println("\n--- Using NewSorter with sort.Sort ---")
	sorterData := []string{"File2.txt", "file10.txt", "FILE1.txt", "item2", "Item10", "ITEM1"}
	fmt.Printf("Before: %v\n", sorterData)

	// Create a case-insensitive sorter
	sorter := ansort.NewSorter(sorterData, ansort.WithCaseInsensitive())
	sort.Sort(sorter)
	fmt.Printf("After (case-insensitive with sort.Sort): %v\n", sorterData)

	// Create a case-sensitive sorter
	sorterData2 := []string{"File2.txt", "file10.txt", "FILE1.txt", "item2", "Item10", "ITEM1"}
	fmt.Printf("Before: %v\n", sorterData2)
	sorter2 := ansort.NewSorter(sorterData2) // default is case-sensitive
	sort.Sort(sorter2)
	fmt.Printf("After (case-sensitive with sort.Sort): %v\n", sorterData2)
}
