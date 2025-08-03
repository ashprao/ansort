package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/ashprao/ansort"
)

func main() {
	fmt.Println("=== Validation Functions Demo ===")

	// 1. Successful validation
	fmt.Println("1. Successful validation:")
	data := []string{"file10.txt", "file2.txt", "file1.txt"}
	err := ansort.SortStringsValidated(data, ansort.WithCaseInsensitive())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sorted data: %v\n\n", data)

	// 2. Validation failure - nil input
	fmt.Println("2. Validation failure - nil input:")
	err = ansort.SortStringsValidated(nil)
	if err != nil {
		var validationErr *ansort.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("ValidationError - Field: %s, Message: %s\n\n",
				validationErr.Field, validationErr.Message)
		}
	}

	// 3. External sort key validation success
	fmt.Println("3. External sort key validation success:")
	key, err := ansort.ToNaturalSortKeyValidated("file10.txt",
		ansort.WithMaxNumericLength(5), ansort.WithExternalCaseInsensitive())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated sort key: %s\n\n", key)

	// 4. External sort key validation failure - invalid padding
	fmt.Println("4. External sort key validation failure - invalid padding:")
	_, err = ansort.ToNaturalSortKeyValidated("file10.txt", ansort.WithMaxNumericLength(0))
	if err != nil {
		var validationErr *ansort.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("ValidationError - Field: %s, Message: %s\n\n",
				validationErr.Field, validationErr.Message)
		}
	}

	// 5. External sort key validation failure - excessive padding
	fmt.Println("5. External sort key validation failure - excessive padding:")
	_, err = ansort.ToNaturalSortKeyValidated("file10.txt", ansort.WithMaxNumericLength(100))
	if err != nil {
		var validationErr *ansort.ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("ValidationError - Field: %s, Message: %s\n\n",
				validationErr.Field, validationErr.Message)
		}
	}

	// 6. Comparison: Validated vs Non-validated
	fmt.Println("6. Comparison: Validated vs Non-validated functions:")
	input := "Item20.txt"

	// Non-validated (convenience)
	keyNonValidated := ansort.ToNaturalSortKey(input, ansort.WithExternalCaseInsensitive())

	// Validated
	keyValidated, err := ansort.ToNaturalSortKeyValidated(input, ansort.WithExternalCaseInsensitive())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Non-validated: %s\n", keyNonValidated)
	fmt.Printf("Validated:     %s\n", keyValidated)
	fmt.Printf("Results match: %t\n", keyNonValidated == keyValidated)
}
