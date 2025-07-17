package utils

// Dedupe removes duplicates from a slice of any comparable type.
func Dedupe[T comparable](input []T) []T {
	seen := make(map[T]bool)

	var result []T

	for _, val := range input {
		if !seen[val] {
			seen[val] = true

			result = append(result, val)
		}
	}

	return result
}
