package main

import (
	"fmt"
	"sort"

	"github.com/ashprao/ansort"
)

func main() {
	fmt.Println("=== External Sort Key Generation Example ===")

	// Example 1: Basic sort key generation
	fmt.Println("\n1. Basic Sort Key Generation:")
	items := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
	fmt.Printf("Original: %v\n", items)

	for _, item := range items {
		key := ansort.ToNaturalSortKey(item)
		fmt.Printf("  %s -> %s\n", item, key)
	}

	// Example 2: Demonstrating lexicographic sorting with generated keys
	fmt.Println("\n2. External System Simulation:")
	fmt.Println("Simulating how Elasticsearch or database would sort using generated keys...")

	type document struct {
		ID      string
		SortKey string
	}

	// Create documents with sort keys
	var docs []document
	for _, item := range items {
		docs = append(docs, document{
			ID:      item,
			SortKey: ansort.ToNaturalSortKey(item),
		})
	}

	// Sort by sort keys (how external system would sort)
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].SortKey < docs[j].SortKey
	})

	fmt.Println("External system sorted order:")
	for _, doc := range docs {
		fmt.Printf("  ID: %s, SortKey: %s\n", doc.ID, doc.SortKey)
	}

	// Example 3: Verify consistency with ansort natural sorting
	fmt.Println("\n3. Consistency Verification:")

	// Natural sort the original items
	naturalSorted := make([]string, len(items))
	copy(naturalSorted, items)
	ansort.SortStrings(naturalSorted)

	// Extract the externally sorted IDs
	externallySorted := make([]string, len(docs))
	for i, doc := range docs {
		externallySorted[i] = doc.ID
	}

	fmt.Printf("Natural sort result:  %v\n", naturalSorted)
	fmt.Printf("External sort result: %v\n", externallySorted)

	// Check if they match
	match := true
	if len(naturalSorted) == len(externallySorted) {
		for i := 0; i < len(naturalSorted); i++ {
			if naturalSorted[i] != externallySorted[i] {
				match = false
				break
			}
		}
	} else {
		match = false
	}

	if match {
		fmt.Println("✅ Results match! External sort keys maintain natural order.")
	} else {
		fmt.Println("❌ Results don't match! There's an issue with sort key generation.")
	}

	// Example 4: Custom configuration options
	fmt.Println("\n4. Custom Configuration Options:")

	// Custom padding length
	fmt.Println("Custom padding length (3 digits):")
	for _, item := range []string{"item1", "item10", "item100"} {
		key := ansort.ToNaturalSortKey(item, ansort.WithMaxNumericLength(3))
		fmt.Printf("  %s -> %s\n", item, key)
	}

	// Case-insensitive keys
	fmt.Println("\nCase-insensitive keys:")
	mixedCase := []string{"File10.TXT", "file2.txt", "FILE1.txt"}
	for _, item := range mixedCase {
		key := ansort.ToNaturalSortKey(item, ansort.WithExternalCaseInsensitive())
		fmt.Printf("  %s -> %s\n", item, key)
	}

	// Example 5: Version strings (common use case)
	fmt.Println("\n5. Version String Example:")
	versions := []string{"v1.10.2", "v1.2.10", "v1.2.3", "v2.1.0"}
	fmt.Printf("Versions: %v\n", versions)

	fmt.Println("Generated sort keys:")
	for _, version := range versions {
		key := ansort.ToNaturalSortKey(version)
		fmt.Printf("  %s -> %s\n", version, key)
	}

	fmt.Println("\n=== Use Cases ===")
	fmt.Println("📊 Elasticsearch: Store sort keys in a 'keyword' field for efficient sorting")
	fmt.Println("🗄️  Databases: Use sort keys in ORDER BY clauses for natural ordering")
	fmt.Println("🔍 Search engines: Enable natural ordering in paginated results")
	fmt.Println("📈 Analytics: Maintain natural order in aggregations and reports")
}
