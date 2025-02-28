package mermaid

import (
	"fmt"
	"slices"
	"strings"
)

// Lexer splits input into tokens
type Lexer struct {
	input string

	tokens []Token

	pos    uint
	line   uint
	column uint
}

// NewLexer creates a new Lexer instance
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
		tokens: make([]Token, 0),
	}
}

// verifyAndSkipHeader verifies header format and returns header info
func (l *Lexer) verifyAndSkipHeader() (*HeaderInfo, error) {
	info := &HeaderInfo{}

	l.skipWhitespace()

	// Parse yaml frontmatter and title if present
	if l.consume("---\n") {
		l.skipWhitespace()
		if l.consume("title: ") {
			titleStart := l.pos
			for l.pos < ulen(l.input) && l.peek() != '\n' {
				l.advance()
			}
			info.Title = l.input[titleStart:l.pos]
			l.advance() // consume newline

			if !l.consume("---\n") {
				return nil, fmt.Errorf("missing YAML frontmatter closing")
			}
		}
	}

	// Verify diagram type
	if l.consume(string(HeaderTypeStateDiagramV1)) {
		info.Type = HeaderTypeStateDiagramV1
	} else if l.consume(string(HeaderTypeStateDiagramV2)) {
		info.Type = HeaderTypeStateDiagramV2
	} else {
		return nil, fmt.Errorf("invalid or missing diagram type declaration")
	}

	// Consume rest of the line
	for l.pos < ulen(l.input) && l.peek() != '\n' {
		l.advance()
	}
	l.advance() // consume newline

	return info, nil
}

// Tokenize processes the input string and returns a list of tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	_, err := l.verifyAndSkipHeader()
	if err != nil {
		return nil, fmt.Errorf("parsing diagram header: %w", err)
	}

	for l.pos < ulen(l.input) {
		ch := l.peek()

		keyword := ""

		word, found := l.peekWord()

		// check for known keywords
		if found && slices.Contains(typedSlice(KeywordValues(), kw.ToStringFunc()), word) {
			keyword = word
		}

		switch {
		case ch == '\n': // add \n as EOL
			l.parseEOL()

		case isWhitespace(ch):
			l.skipWhitespace()

		case keyword != "":
			err := l.parseKeyword(keyword)
			if err != nil {
				return nil, err
			}

		case ch == KeywordColon[0]:
			l.addToken(TokenTypeKeyword, KeywordColon.String())
			l.advance()

		case l.consume(KeywordArrow.String()):
			l.addToken(TokenTypeKeyword, KeywordArrow.String())

		case ch == '[':
			_, found := l.consumeTill(false, `*]`) // disrespect new line
			if !found {
				return nil, fmt.Errorf("weird [ on line (%d)", l.line)
			}

			l.addToken(TokenTypeIdentifier, KeywordTerminalState.String())
			l.advanceN(2) // skip '*]'

		case ch == '"':
			l.advance() // skip '"'

			res, found := l.consumeTill(false, `"`) // disrespect new line
			if !found {
				return nil, fmt.Errorf("missing closing quote on line (%d) after opening quote on position (%d). multiline strings are not supported yet", l.line, l.pos)
			}

			l.addToken(TokenTypeText, res)
			l.advance() // skip '"'

		default:
			l.parseText(ch)
		}
	}

	l.addToken(TokenTypeEOF, "")

	// delete duplicate EOL tokens
	for i := uint(0); i < ulen(l.tokens)-1; i++ {
		// both are EOL
		if (l.tokens[i].Type == TokenTypeEOL && l.tokens[i+1].Type == TokenTypeEOL) &&
			// and both share same column and line
			(l.tokens[i].Column == l.tokens[i+1].Column &&
				l.tokens[i].Line == l.tokens[i+1].Line) {

			l.tokens = append(l.tokens[:i], l.tokens[i+1:]...)
			i--

		}
	}

	return l.tokens, nil
}

func (l *Lexer) parseText(ch byte) {
	// Handle text
	if isAlphaNumeric(ch) {
		start := l.pos
		for l.pos < ulen(l.input) && (isAlphaNumeric(l.peek()) || l.peek() == '_') {
			l.advance()
		}
		l.addToken(TokenTypeIdentifier, l.input[start:l.pos])
	} else {
		fmt.Printf("unexpected character (%s) at line (%d), peek 10 (%s)\n", string(ch), l.line, l.peekN(10))
		l.advance()
	}
}

func (l *Lexer) parseKeyword(keyword string) error {
	switch keyword {
	case KeywordBlockStart.String():
		l.addToken(TokenTypeBlockStart, KeywordBlockStart.String())
		l.advanceN(len(keyword))

	case KeywordBlockEnd.String():
		l.addToken(TokenTypeBlockEnd, KeywordBlockEnd.String())
		l.advanceN(len(keyword))

	case KeywordColon.String():
		l.advance() // skip ":"

		res, found := l.consumeTillLineBreak()
		if !found {
			return fmt.Errorf("something is wrong with ':' at line (%d)", l.line)
		}

		l.addToken(TokenTypeKeyword, KeywordColon.String())
		l.addToken(TokenTypeText, strings.TrimSpace(res))
		l.addToken(TokenTypeEOL, "\n")

	case KeywordArrow.String():
		l.addToken(TokenTypeKeyword, keyword)
		l.advanceN(len(keyword))

	case KeywordComment.String():
		l.advanceN(2) // skip "%%"

		res, found := l.consumeTillLineBreak()
		if !found {
			return fmt.Errorf("something is wrong with comment at line (%d)", l.line)
		}

		l.addToken(TokenTypeText, strings.TrimSpace(res))
		l.addToken(TokenTypeEOL, "\n")

	default:
		l.addToken(TokenTypeKeyword, keyword)
		l.advanceN(len(keyword))
	}
	return nil

}

func (l *Lexer) parseEOL() {
	l.addToken(TokenTypeEOL, "\n")
	l.advance()
}
