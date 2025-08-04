package main

import (
	"fmt"
	"os"

	"github.com/ashprao/ansort"
)

func main() {
	// Change to the parent directory to access the ansort package
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}

	// Test the comparison
	original := ansort.Compare("item001", "item1")
	fmt.Printf("Original Compare('item001', 'item1'): %d\n", original)

	// Test individual token parsing to see what's happening
	fmt.Println("This helps us debug the comparison logic")
}
