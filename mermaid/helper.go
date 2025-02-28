package mermaid

import (
	"strings"
	"unicode"
)

type ulenType interface {
	string | []Token
}

// ulen returns the length of a string as a uint
func ulen[T ulenType](s T) uint {
	return uint(len(s))
}

// isWhitespace checks if a character is whitespace
func isWhitespace(ch byte) bool {
	return unicode.IsSpace(rune(ch))
}

// isAlphaNumeric checks if a character is alphanumeric
func isAlphaNumeric(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch))
}

// isWordChar checks if a character is part of a "word" (alphanumeric or part of a symbol/token)
func isWordChar(ch byte) bool {
	// Consider alphanumeric characters and common symbols as valid word characters
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || strings.ContainsRune("_-:<>%{}", rune(ch))
}
