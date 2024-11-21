package txt

import (
	"strings"
)

// CapitalizeWords capitalizes the first letter of each word in a string.
func CapitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		words[i] = FirstToUpper(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}
