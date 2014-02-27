package common

import (
	"strconv"
)

func HashCode(in string) int32 {
	// Initialize output
	var hash int32
	// Empty string has a hashcode of 0
	if len(in) == 0 {
		return hash
	}
	// Convert string into slice of bytes
	b := []byte(in)
	// Build hash
	for i := range b {
		char := b[i]
		hash = ((hash << 5) - hash) + int32(char)
	}
	return hash
}
func HashString(in string) string {
	code := HashCode(in)
	return strconv.Itoa(int(code))
}
