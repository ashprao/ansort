package main

import "fmt"

func main() {
	s := fmt.Sprintf("very_long_filename_with_lots_of_characters_%d.txt", 45)
	fmt.Printf("String: %s\nLength: %d\n", s, len(s))
}
