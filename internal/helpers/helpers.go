// Package helpers provides utility functions for various tasks.
package helpers

import "fmt"

func IncrementCode(code string) string {
	var prefix string
	var number int
	fmt.Sscanf(code, "%2s-%d", &prefix, &number)
	number++ // Increment the numeric part

	return fmt.Sprintf("%s-%02d", prefix, number) // Format it back to the same structure
}
