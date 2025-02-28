package mermaid

import "strings"

// peek returns the next character without advancing
func (l *Lexer) peek() byte {
	if l.pos >= ulen(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// advance moves to the next character
func (l *Lexer) advance() {
	if l.pos < ulen(l.input) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

// advanceN advances position by n characters
func (l *Lexer) advanceN(n int) {
	for i := 0; i < n && l.pos < ulen(l.input); i++ {
		l.advance()
	}
}

// skipWhitespace skips over any whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.pos < ulen(l.input) && isWhitespace(l.peek()) {
		l.advance()
	}
}

// addToken adds a new token to the token list
func (l *Lexer) addToken(typ TokenType, value string) {
	l.tokens = append(l.tokens, Token{
		Type:   typ,
		Value:  value,
		Line:   l.line,
		Column: l.column,
	})
}

// peekN returns the next n characters without advancing position
func (l *Lexer) peekN(n uint) string {
	if l.pos+n > ulen(l.input) {
		return ""
	}
	return l.input[l.pos : l.pos+n]
}

// peekTill returns the string up to the first occurrence of any of the matches from the current position.
// If this does respect newline, it will stop looking for a match if it's not on the current line
// It does not advance the position.
// Matches are exclusive.
// It will return the whole input (or till newline) if no match is found
func (l *Lexer) peekTill(respectNewLine bool, matches ...string) (string, bool) {
	// Store the initial position to restore after peeking
	startPos := l.pos

	// Loop through the input until the end or a newline is encountered
	for l.pos < ulen(l.input) {
		if respectNewLine && l.peek() == '\n' {
			break
		}

		// Check each match string to see if it starts at the current position
		for _, match := range matches {
			if strings.HasPrefix(l.input[l.pos:], match) {
				// Return the substring from the start position to the match position
				result := l.input[startPos:l.pos]
				// Restore the cursor position
				l.pos = startPos
				return result, true
			}
		}
		// Temporarily advance to check for matches
		l.advance()
	}

	// Save last position seeked to
	lastPos := l.pos
	// If no match is found, restore the cursor position
	l.pos = startPos

	// return the whole string if no match was found
	return l.input[startPos:lastPos], false
}

// consumeTillLineBreak returns the content up to the first newline (\n) or end of input,
// and advances the position to the character after the newline.
func (l *Lexer) consumeTillLineBreak() (string, bool) {
	result, _ := l.peekTill(true)
	l.advanceN(len(result) + 1) // + newline
	return result, true
}

// consume advances position by n characters if they match the expected string
func (l *Lexer) consume(expected string) bool {
	if l.peekN(ulen(expected)) == expected {
		l.advanceN(len(expected))

		return true
	}
	return false
}

// Generic function to convert a slice of one type to another
func typedSlice[T any, U any](input []T, convertFunc func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = convertFunc(v)
	}
	return result
}

// consumeTill consumes the input up to the first occurrence of any of the matches
// and advances the position. If no match is found, returns false.
// consumeTill advances the position and returns the string up to the first occurrence of any of the matches.
// If no match is found, it returns false.
// This does respect newline, so it will stop looking for a match if it's not on the current line
// Matches are exclusive.
func (l *Lexer) consumeTill(respectNewLine bool, matches ...string) (string, bool) {
	// Use peekTill to find the substring and determine if a match exists
	result, found := l.peekTill(respectNewLine, matches...)
	if found {
		// Advance the position by the length of the result
		l.advanceN(len(result))
	}
	return result, found
}

// peekWord returns the word or token under the current cursor without advancing the position.
// A word is defined as a sequence of valid "word characters."
func (l *Lexer) peekWord() (string, bool) {
	startPos := l.pos

	// Skip leading whitespace
	for startPos < ulen(l.input) && isWhitespace(l.input[startPos]) {
		startPos++
	}

	// If we hit the end of input, return false
	if startPos >= ulen(l.input) {
		return "", false
	}

	endPos := startPos

	// Find the end of the word (sequence of valid word characters)
	for endPos < ulen(l.input) && isWordChar(l.input[endPos]) {
		endPos++
	}

	// Return the word without modifying the cursor
	return l.input[startPos:endPos], true
}
