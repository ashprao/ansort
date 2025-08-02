package main

import (
	"fmt"

	"github.com/ashprao/ansort"
)

func main() {
	// Example 1: Basic sorting (case-sensitive by default)
	fmt.Println("=== Basic Alphanumeric Sorting (Case-Sensitive) ===")
	data1 := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
	fmt.Printf("Before: %v\n", data1)
	ansort.SortStrings(data1)
	fmt.Printf("After:  %v\n", data1)

	// Example 2: Case-insensitive sorting
	fmt.Println("\n=== Case-Insensitive Sorting ===")
	data2 := []string{"File10.txt", "file2.txt", "FILE1.txt", "file20.txt"}
	fmt.Printf("Before: %v\n", data2)
	ansort.SortStrings(data2, ansort.WithCaseInsensitive())
	fmt.Printf("After:  %v\n", data2)

	// Example 3: Version-like strings
	fmt.Println("\n=== Version Strings ===")
	data3 := []string{"v1.10.0", "v1.2.0", "v1.20.0", "v1.1.0"}
	fmt.Printf("Before: %v\n", data3)
	ansort.SortStrings(data3)
	fmt.Printf("After:  %v\n", data3)

	// Example 4: Mixed content
	fmt.Println("\n=== Mixed Content ===")
	data4 := []string{"item10", "item2", "item1", "item100", "item20"}
	fmt.Printf("Before: %v\n", data4)
	ansort.SortStrings(data4)
	fmt.Printf("After:  %v\n", data4)

	// Example 5: String comparison (case-sensitive by default)
	fmt.Println("\n=== String Comparison (Case-Sensitive) ===")
	a, b := "file1.txt", "file10.txt"
	result := ansort.Compare(a, b)
	fmt.Printf("Compare(%q, %q) = %d\n", a, b, result)

	a, b = "file10.txt", "file2.txt"
	result = ansort.Compare(a, b)
	fmt.Printf("Compare(%q, %q) = %d\n", a, b, result)

	// Example 6: Case-insensitive comparison
	fmt.Println("\n=== String Comparison (Case-Insensitive) ===")
	a, b = "File1.txt", "file10.txt"
	result = ansort.Compare(a, b, ansort.WithCaseInsensitive())
	fmt.Printf("Compare(%q, %q) with case-insensitive = %d\n", a, b, result)

	a, b = "FILE1.txt", "file1.txt"
	result = ansort.Compare(a, b, ansort.WithCaseInsensitive())
	fmt.Printf("Compare(%q, %q) with case-insensitive = %d\n", a, b, result)
}
