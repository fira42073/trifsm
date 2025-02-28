package mermaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name       string
		inputState string
		result     []Token
	}{
		{
			name: "Simple without frontmatter",
			inputState: `stateDiagram-v2
	state Still`,
			result: []Token{
				{Type: TokenTypeKeyword, Line: 2, Column: 2, Value: "state"},     // 0
				{Type: TokenTypeIdentifier, Line: 2, Column: 13, Value: "Still"}, // 1
				{Type: TokenTypeEOF, Line: 2, Column: 13, Value: ""},             // 2
			},
		},
		{
			name: "Simple with frontmatter",
			inputState: `---
title: Simple sample
---
stateDiagram-v2
	state Still
	state "State of crash" as Crash
	Moving: Description of moving state

    [*] --> Still
    Still --> [*]

    Still --> Moving
    Moving --> Still
    Moving --> Crash
    Crash --> [*]
	note right of Still: This is a note
`,
			result: []Token{
				{Type: TokenTypeKeyword, Line: 5, Column: 2, Value: "state"},                            // 0
				{Type: TokenTypeIdentifier, Line: 5, Column: 13, Value: "Still"},                        // 1
				{Type: TokenTypeEOL, Line: 5, Column: 13, Value: "\n"},                                  // 2
				{Type: TokenTypeKeyword, Line: 6, Column: 2, Value: "state"},                            // 3
				{Type: TokenTypeText, Line: 6, Column: 37, Value: "State of crash"},                     // 4
				{Type: TokenTypeKeyword, Line: 6, Column: 39, Value: KeywordAs.String()},                // 5
				{Type: TokenTypeIdentifier, Line: 6, Column: 47, Value: "Crash"},                        // 6
				{Type: TokenTypeEOL, Line: 6, Column: 47, Value: "\n"},                                  // 7
				{Type: TokenTypeIdentifier, Line: 7, Column: 8, Value: "Moving"},                        // 8
				{Type: TokenTypeKeyword, Line: 8, Column: 1, Value: KeywordColon.String()},              // 9
				{Type: TokenTypeText, Line: 8, Column: 1, Value: "Description of moving state"},         // 10
				{Type: TokenTypeEOL, Line: 8, Column: 1, Value: "\n"},                                   // 11
				{Type: TokenTypeIdentifier, Line: 9, Column: 7, Value: KeywordTerminalState.String()},   // 12
				{Type: TokenTypeKeyword, Line: 9, Column: 10, Value: KeywordArrow.String()},             // 13
				{Type: TokenTypeIdentifier, Line: 9, Column: 19, Value: "Still"},                        // 14
				{Type: TokenTypeEOL, Line: 9, Column: 19, Value: "\n"},                                  // 15
				{Type: TokenTypeIdentifier, Line: 10, Column: 10, Value: "Still"},                       // 16
				{Type: TokenTypeKeyword, Line: 10, Column: 11, Value: KeywordArrow.String()},            // 17
				{Type: TokenTypeIdentifier, Line: 10, Column: 17, Value: KeywordTerminalState.String()}, // 18
				{Type: TokenTypeEOL, Line: 10, Column: 19, Value: "\n"},                                 // 19
				{Type: TokenTypeEOL, Line: 11, Column: 1, Value: "\n"},                                  // 20
				{Type: TokenTypeIdentifier, Line: 12, Column: 10, Value: "Still"},                       // 21
				{Type: TokenTypeKeyword, Line: 12, Column: 11, Value: KeywordArrow.String()},            // 22
				{Type: TokenTypeIdentifier, Line: 12, Column: 21, Value: "Moving"},                      // 23
				{Type: TokenTypeEOL, Line: 12, Column: 21, Value: "\n"},                                 // 24
				{Type: TokenTypeIdentifier, Line: 13, Column: 11, Value: "Moving"},                      // 25
				{Type: TokenTypeKeyword, Line: 13, Column: 12, Value: KeywordArrow.String()},            // 26
				{Type: TokenTypeIdentifier, Line: 13, Column: 21, Value: "Still"},                       // 27
				{Type: TokenTypeEOL, Line: 13, Column: 21, Value: "\n"},                                 // 28
				{Type: TokenTypeIdentifier, Line: 14, Column: 11, Value: "Moving"},                      // 29
				{Type: TokenTypeKeyword, Line: 14, Column: 12, Value: KeywordArrow.String()},            // 30
				{Type: TokenTypeIdentifier, Line: 14, Column: 21, Value: "Crash"},                       // 31
				{Type: TokenTypeEOL, Line: 14, Column: 21, Value: "\n"},                                 // 32
				{Type: TokenTypeIdentifier, Line: 15, Column: 10, Value: "Crash"},                       // 33
				{Type: TokenTypeKeyword, Line: 15, Column: 11, Value: KeywordArrow.String()},            // 34
				{Type: TokenTypeIdentifier, Line: 15, Column: 17, Value: KeywordTerminalState.String()}, // 35
				{Type: TokenTypeEOL, Line: 15, Column: 19, Value: "\n"},                                 // 36
				{Type: TokenTypeKeyword, Line: 16, Column: 2, Value: "note"},                            // 37
				{Type: TokenTypeKeyword, Line: 16, Column: 7, Value: "right"},                           // 38
				{Type: TokenTypeKeyword, Line: 16, Column: 13, Value: "of"},                             // 39
				{Type: TokenTypeIdentifier, Line: 16, Column: 21, Value: "Still"},                       // 40
				{Type: TokenTypeKeyword, Line: 17, Column: 1, Value: KeywordColon.String()},             // 41
				{Type: TokenTypeText, Line: 17, Column: 1, Value: "This is a note"},                     // 42
				{Type: TokenTypeEOL, Line: 17, Column: 1, Value: "\n"},                                  // 43
				{Type: TokenTypeEOF, Line: 17, Column: 1, Value: ""},                                    // 44
			},
		},
		{
			name: "All features",
			inputState: `---
title: Simple sample (optional header)
---
stateDiagram-v2
    direction LR

    StateB
    StateC: C declr

    [*] --> StateA

    state "State A" as StateA
    StateA --> StateB: Event1
    StateB --> StateC: Event2
	StateC --> [*]

    note right of StateA: This is a note right of StateA
    note left of StateB: This is a note left of StateB

    %% comment here for StateA
    StateA --> Choice1: Decision Point
    StateB --> <<fork>> ForkState
    ForkState --> StateD
    ForkState --> StateE
    StateD --> <<join>> JoinState
    StateE --> <<join>> JoinState
    JoinState --> StateF

    StateF --> <<choice>> Choice1
    Choice1 --> StateG: Option 1
    Choice1 --> StateH: Option 2
    StateG --> ConcurrentState
    StateH --> ConcurrentState

    state ConcurrentState {
        direction LR
        [*] --> SubState1
        SubState1 --> SubState2
        --
        SubState2 --> [*]
    }

    ConcurrentState --> FinalState: Exit
    FinalState --> [*]

    state "Final State" as FinalState`,
			result: []Token{
				{Type: TokenTypeKeyword, Line: 5, Column: 5, Value: "direction"},                        // 0
				{Type: TokenTypeIdentifier, Line: 5, Column: 17, Value: "LR"},                           // 1
				{Type: TokenTypeEOL, Line: 5, Column: 17, Value: "\n"},                                  // 2
				{Type: TokenTypeEOL, Line: 6, Column: 1, Value: "\n"},                                   // 3
				{Type: TokenTypeIdentifier, Line: 7, Column: 11, Value: "StateB"},                       // 4
				{Type: TokenTypeEOL, Line: 7, Column: 11, Value: "\n"},                                  // 5
				{Type: TokenTypeIdentifier, Line: 8, Column: 11, Value: "StateC"},                       // 6
				{Type: TokenTypeKeyword, Line: 9, Column: 1, Value: KeywordColon.String()},              // 7
				{Type: TokenTypeText, Line: 9, Column: 1, Value: "C declr"},                             // 8
				{Type: TokenTypeEOL, Line: 9, Column: 1, Value: "\n"},                                   // 9
				{Type: TokenTypeIdentifier, Line: 10, Column: 7, Value: KeywordTerminalState.String()},  // 10
				{Type: TokenTypeKeyword, Line: 10, Column: 10, Value: KeywordArrow.String()},            // 11
				{Type: TokenTypeIdentifier, Line: 10, Column: 20, Value: "StateA"},                      // 12
				{Type: TokenTypeEOL, Line: 10, Column: 20, Value: "\n"},                                 // 13
				{Type: TokenTypeEOL, Line: 11, Column: 1, Value: "\n"},                                  // 14
				{Type: TokenTypeKeyword, Line: 12, Column: 5, Value: "state"},                           // 15
				{Type: TokenTypeText, Line: 12, Column: 26, Value: "State A"},                           // 16
				{Type: TokenTypeKeyword, Line: 12, Column: 28, Value: KeywordAs.String()},               // 17
				{Type: TokenTypeIdentifier, Line: 12, Column: 37, Value: "StateA"},                      // 18
				{Type: TokenTypeEOL, Line: 12, Column: 37, Value: "\n"},                                 // 19
				{Type: TokenTypeIdentifier, Line: 13, Column: 11, Value: "StateA"},                      // 20
				{Type: TokenTypeKeyword, Line: 13, Column: 12, Value: KeywordArrow.String()},            // 21
				{Type: TokenTypeIdentifier, Line: 13, Column: 22, Value: "StateB"},                      // 22
				{Type: TokenTypeKeyword, Line: 14, Column: 1, Value: KeywordColon.String()},             // 23
				{Type: TokenTypeText, Line: 14, Column: 1, Value: "Event1"},                             // 24
				{Type: TokenTypeEOL, Line: 14, Column: 1, Value: "\n"},                                  // 25
				{Type: TokenTypeIdentifier, Line: 14, Column: 11, Value: "StateB"},                      // 26
				{Type: TokenTypeKeyword, Line: 14, Column: 12, Value: KeywordArrow.String()},            // 27
				{Type: TokenTypeIdentifier, Line: 14, Column: 22, Value: "StateC"},                      // 28
				{Type: TokenTypeKeyword, Line: 15, Column: 1, Value: KeywordColon.String()},             // 29
				{Type: TokenTypeText, Line: 15, Column: 1, Value: "Event2"},                             // 30
				{Type: TokenTypeEOL, Line: 15, Column: 1, Value: "\n"},                                  // 31
				{Type: TokenTypeIdentifier, Line: 15, Column: 8, Value: "StateC"},                       // 32
				{Type: TokenTypeKeyword, Line: 15, Column: 9, Value: KeywordArrow.String()},             // 33
				{Type: TokenTypeIdentifier, Line: 15, Column: 15, Value: KeywordTerminalState.String()}, // 34
				{Type: TokenTypeEOL, Line: 15, Column: 17, Value: "\n"},                                 // 35
				{Type: TokenTypeEOL, Line: 16, Column: 1, Value: "\n"},                                  // 36
				{Type: TokenTypeKeyword, Line: 17, Column: 5, Value: "note"},                            // 37
				{Type: TokenTypeKeyword, Line: 17, Column: 10, Value: "right"},                          // 38
				{Type: TokenTypeKeyword, Line: 17, Column: 16, Value: "of"},                             // 39
				{Type: TokenTypeIdentifier, Line: 17, Column: 25, Value: "StateA"},                      // 40
				{Type: TokenTypeKeyword, Line: 18, Column: 1, Value: KeywordColon.String()},             // 41
				{Type: TokenTypeText, Line: 18, Column: 1, Value: "This is a note right of StateA"},     // 42
				{Type: TokenTypeEOL, Line: 18, Column: 1, Value: "\n"},                                  // 43
				{Type: TokenTypeKeyword, Line: 18, Column: 5, Value: "note"},                            // 44
				{Type: TokenTypeKeyword, Line: 18, Column: 10, Value: "left"},                           // 45
				{Type: TokenTypeKeyword, Line: 18, Column: 15, Value: "of"},                             // 46
				{Type: TokenTypeIdentifier, Line: 18, Column: 24, Value: "StateB"},                      // 47
				{Type: TokenTypeKeyword, Line: 19, Column: 1, Value: KeywordColon.String()},             // 48
				{Type: TokenTypeText, Line: 19, Column: 1, Value: "This is a note left of StateB"},      // 49
				{Type: TokenTypeEOL, Line: 19, Column: 1, Value: "\n"},                                  // 50
				{Type: TokenTypeText, Line: 21, Column: 1, Value: "comment here for StateA"},            // 51
				{Type: TokenTypeEOL, Line: 21, Column: 1, Value: "\n"},                                  // 52
				{Type: TokenTypeIdentifier, Line: 21, Column: 11, Value: "StateA"},                      // 53
				{Type: TokenTypeKeyword, Line: 21, Column: 12, Value: KeywordArrow.String()},            // 54
				{Type: TokenTypeIdentifier, Line: 21, Column: 23, Value: "Choice1"},                     // 55
				{Type: TokenTypeKeyword, Line: 22, Column: 1, Value: KeywordColon.String()},             // 56
				{Type: TokenTypeText, Line: 22, Column: 1, Value: "Decision Point"},                     // 57
				{Type: TokenTypeEOL, Line: 22, Column: 1, Value: "\n"},                                  // 58
				{Type: TokenTypeIdentifier, Line: 22, Column: 11, Value: "StateB"},                      // 59
				{Type: TokenTypeKeyword, Line: 22, Column: 12, Value: KeywordArrow.String()},            // 60
				{Type: TokenTypeKeyword, Line: 22, Column: 16, Value: KeywordFork.String()},             // 61
				{Type: TokenTypeIdentifier, Line: 22, Column: 34, Value: "ForkState"},                   // 62
				{Type: TokenTypeEOL, Line: 22, Column: 34, Value: "\n"},                                 // 63
				{Type: TokenTypeIdentifier, Line: 23, Column: 14, Value: "ForkState"},                   // 64
				{Type: TokenTypeKeyword, Line: 23, Column: 15, Value: KeywordArrow.String()},            // 65
				{Type: TokenTypeIdentifier, Line: 23, Column: 25, Value: "StateD"},                      // 66
				{Type: TokenTypeEOL, Line: 23, Column: 25, Value: "\n"},                                 // 67
				{Type: TokenTypeIdentifier, Line: 24, Column: 14, Value: "ForkState"},                   // 68
				{Type: TokenTypeKeyword, Line: 24, Column: 15, Value: KeywordArrow.String()},            // 69
				{Type: TokenTypeIdentifier, Line: 24, Column: 25, Value: "StateE"},                      // 70
				{Type: TokenTypeEOL, Line: 24, Column: 25, Value: "\n"},                                 // 71
				{Type: TokenTypeIdentifier, Line: 25, Column: 11, Value: "StateD"},                      // 72
				{Type: TokenTypeKeyword, Line: 25, Column: 12, Value: KeywordArrow.String()},            // 73
				{Type: TokenTypeKeyword, Line: 25, Column: 16, Value: KeywordJoin.String()},             // 74
				{Type: TokenTypeIdentifier, Line: 25, Column: 34, Value: "JoinState"},                   // 75
				{Type: TokenTypeEOL, Line: 25, Column: 34, Value: "\n"},                                 // 76
				{Type: TokenTypeIdentifier, Line: 26, Column: 11, Value: "StateE"},                      // 77
				{Type: TokenTypeKeyword, Line: 26, Column: 12, Value: KeywordArrow.String()},            // 78
				{Type: TokenTypeKeyword, Line: 26, Column: 16, Value: KeywordJoin.String()},             // 79
				{Type: TokenTypeIdentifier, Line: 26, Column: 34, Value: "JoinState"},                   // 80
				{Type: TokenTypeEOL, Line: 26, Column: 34, Value: "\n"},                                 // 81
				{Type: TokenTypeIdentifier, Line: 27, Column: 14, Value: "JoinState"},                   // 82
				{Type: TokenTypeKeyword, Line: 27, Column: 15, Value: KeywordArrow.String()},            // 83
				{Type: TokenTypeIdentifier, Line: 27, Column: 25, Value: "StateF"},                      // 84
				{Type: TokenTypeEOL, Line: 27, Column: 25, Value: "\n"},                                 // 85
				{Type: TokenTypeEOL, Line: 28, Column: 1, Value: "\n"},                                  // 86
				{Type: TokenTypeIdentifier, Line: 29, Column: 11, Value: "StateF"},                      // 87
				{Type: TokenTypeKeyword, Line: 29, Column: 12, Value: KeywordArrow.String()},            // 88
				{Type: TokenTypeKeyword, Line: 29, Column: 16, Value: "<<choice>>"},                     // 89
				{Type: TokenTypeIdentifier, Line: 29, Column: 34, Value: "Choice1"},                     // 90
				{Type: TokenTypeEOL, Line: 29, Column: 34, Value: "\n"},                                 // 91
				{Type: TokenTypeIdentifier, Line: 30, Column: 12, Value: "Choice1"},                     // 92
				{Type: TokenTypeKeyword, Line: 30, Column: 13, Value: KeywordArrow.String()},            // 93
				{Type: TokenTypeIdentifier, Line: 30, Column: 23, Value: "StateG"},                      // 94
				{Type: TokenTypeKeyword, Line: 31, Column: 1, Value: KeywordColon.String()},             // 95
				{Type: TokenTypeText, Line: 31, Column: 1, Value: "Option 1"},                           // 96
				{Type: TokenTypeEOL, Line: 31, Column: 1, Value: "\n"},                                  // 97
				{Type: TokenTypeIdentifier, Line: 31, Column: 12, Value: "Choice1"},                     // 98
				{Type: TokenTypeKeyword, Line: 31, Column: 13, Value: KeywordArrow.String()},            // 99
				{Type: TokenTypeIdentifier, Line: 31, Column: 23, Value: "StateH"},                      // 100
				{Type: TokenTypeKeyword, Line: 32, Column: 1, Value: KeywordColon.String()},             // 101
				{Type: TokenTypeText, Line: 32, Column: 1, Value: "Option 2"},                           // 102
				{Type: TokenTypeEOL, Line: 32, Column: 1, Value: "\n"},                                  // 103
				{Type: TokenTypeIdentifier, Line: 32, Column: 11, Value: "StateG"},                      // 104
				{Type: TokenTypeKeyword, Line: 32, Column: 12, Value: KeywordArrow.String()},            // 105
				{Type: TokenTypeIdentifier, Line: 32, Column: 31, Value: "ConcurrentState"},             // 106
				{Type: TokenTypeEOL, Line: 32, Column: 31, Value: "\n"},                                 // 107
				{Type: TokenTypeIdentifier, Line: 33, Column: 11, Value: "StateH"},                      // 108
				{Type: TokenTypeKeyword, Line: 33, Column: 12, Value: KeywordArrow.String()},            // 109
				{Type: TokenTypeIdentifier, Line: 33, Column: 31, Value: "ConcurrentState"},             // 110
				{Type: TokenTypeEOL, Line: 33, Column: 31, Value: "\n"},                                 // 111
				{Type: TokenTypeEOL, Line: 34, Column: 1, Value: "\n"},                                  // 112
				{Type: TokenTypeKeyword, Line: 35, Column: 5, Value: "state"},                           // 113
				{Type: TokenTypeIdentifier, Line: 35, Column: 26, Value: "ConcurrentState"},             // 114
				{Type: TokenTypeBlockStart, Line: 35, Column: 27, Value: "{"},                           // 115
				{Type: TokenTypeEOL, Line: 35, Column: 28, Value: "\n"},                                 // 116
				{Type: TokenTypeKeyword, Line: 36, Column: 9, Value: "direction"},                       // 117
				{Type: TokenTypeIdentifier, Line: 36, Column: 21, Value: "LR"},                          // 118
				{Type: TokenTypeEOL, Line: 36, Column: 21, Value: "\n"},                                 // 119
				{Type: TokenTypeIdentifier, Line: 37, Column: 11, Value: KeywordTerminalState.String()}, // 120
				{Type: TokenTypeKeyword, Line: 37, Column: 14, Value: KeywordArrow.String()},            // 121
				{Type: TokenTypeIdentifier, Line: 37, Column: 27, Value: "SubState1"},                   // 122
				{Type: TokenTypeEOL, Line: 37, Column: 27, Value: "\n"},                                 // 123
				{Type: TokenTypeIdentifier, Line: 38, Column: 18, Value: "SubState1"},                   // 124
				{Type: TokenTypeKeyword, Line: 38, Column: 19, Value: KeywordArrow.String()},            // 125
				{Type: TokenTypeIdentifier, Line: 38, Column: 32, Value: "SubState2"},                   // 126
				{Type: TokenTypeEOL, Line: 38, Column: 32, Value: "\n"},                                 // 127
				{Type: TokenTypeKeyword, Line: 39, Column: 9, Value: KeywordConcurrent.String()},        // 128
				{Type: TokenTypeEOL, Line: 39, Column: 11, Value: "\n"},                                 // 129
				{Type: TokenTypeIdentifier, Line: 40, Column: 18, Value: "SubState2"},                   // 130
				{Type: TokenTypeKeyword, Line: 40, Column: 19, Value: KeywordArrow.String()},            // 131
				{Type: TokenTypeIdentifier, Line: 40, Column: 25, Value: KeywordTerminalState.String()}, // 132
				{Type: TokenTypeEOL, Line: 40, Column: 27, Value: "\n"},                                 // 133
				{Type: TokenTypeBlockEnd, Line: 41, Column: 5, Value: "}"},                              // 134
				{Type: TokenTypeEOL, Line: 41, Column: 6, Value: "\n"},                                  // 135
				{Type: TokenTypeEOL, Line: 42, Column: 1, Value: "\n"},                                  // 136
				{Type: TokenTypeIdentifier, Line: 43, Column: 20, Value: "ConcurrentState"},             // 137
				{Type: TokenTypeKeyword, Line: 43, Column: 21, Value: KeywordArrow.String()},            // 138
				{Type: TokenTypeIdentifier, Line: 43, Column: 35, Value: "FinalState"},                  // 139
				{Type: TokenTypeKeyword, Line: 44, Column: 1, Value: KeywordColon.String()},             // 140
				{Type: TokenTypeText, Line: 44, Column: 1, Value: "Exit"},                               // 141
				{Type: TokenTypeEOL, Line: 44, Column: 1, Value: "\n"},                                  // 142
				{Type: TokenTypeIdentifier, Line: 44, Column: 15, Value: "FinalState"},                  // 143
				{Type: TokenTypeKeyword, Line: 44, Column: 16, Value: KeywordArrow.String()},            // 144
				{Type: TokenTypeIdentifier, Line: 44, Column: 22, Value: KeywordTerminalState.String()}, // 145
				{Type: TokenTypeEOL, Line: 44, Column: 24, Value: "\n"},                                 // 146
				{Type: TokenTypeEOL, Line: 45, Column: 1, Value: "\n"},                                  // 147
				{Type: TokenTypeKeyword, Line: 46, Column: 5, Value: "state"},                           // 148
				{Type: TokenTypeText, Line: 46, Column: 34, Value: "Final State"},                       // 149
				{Type: TokenTypeKeyword, Line: 46, Column: 36, Value: KeywordAs.String()},               // 150
				{Type: TokenTypeIdentifier, Line: 46, Column: 49, Value: "FinalState"},                  // 151
				{Type: TokenTypeEOF, Line: 46, Column: 49, Value: ""},                                   // 152
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create tokenizer and get tokens
			tokenizer := NewLexer(tc.inputState)
			tokens, err := tokenizer.Tokenize()
			if err != nil {
				panic(err)
			}

			for i, token := range tokens {
				assert.Equal(t, tc.result[i].Type, token.Type)
				assert.Equal(t, tc.result[i].Value, token.Value)
				assert.Equal(t, tc.result[i].Line, token.Line)
				assert.Equal(t, tc.result[i].Column, token.Column)
			}
		})
	}
}

func TestPeekTill(t *testing.T) {
	inputText := `hello world! This is a test.
	multistring one`
	tests := []struct {
		name           string
		input          []string
		expected       string
		found          bool
		resetAfter     bool // Reset tokenizer after the test
		respectNewLine bool // Respect line breaks
	}{
		{
			name:     "Match single string",
			input:    []string{"world"},
			expected: "hello ",
			found:    true,
		},
		{
			name:     "Match multiple strings",
			input:    []string{"world", "This"},
			expected: "hello ",
			found:    true,
		},
		{
			name:     "No match",
			input:    []string{"notfound"},
			expected: inputText, // the whole input
			found:    false,
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: inputText,
			found:    false,
		},
		{
			name:       "Position unchanged after peek",
			input:      []string{"world"},
			expected:   "hello ",
			found:      true,
			resetAfter: true,
		},
		{
			name:           "Respect line breaks",
			input:          []string{"multistring"},
			expected:       `hello world! This is a test.`,
			found:          false,
			resetAfter:     true,
			respectNewLine: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewLexer(inputText)

			result, found := tokenizer.peekTill(test.respectNewLine, test.input...)

			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.found, found)

			if test.resetAfter {
				tokenizer.pos = 0 // Reset position for position-unchanged test
			}
		})
	}
}

func TestConsumeTill(t *testing.T) {
	inputText := `hello world! This is a test.
	multistring one`

	tests := []struct {
		name                 string
		input                []string
		expected             string
		found                bool
		positionShouldChange bool
		respectNewLine       bool
	}{
		{
			name:                 "Match single string",
			input:                []string{"world"},
			expected:             "hello ",
			found:                true,
			positionShouldChange: true,
		},
		{
			name:                 "Match multiple strings",
			input:                []string{"!", "This"},
			expected:             "hello world",
			found:                true,
			positionShouldChange: true,
		},
		{
			name:                 "No match",
			input:                []string{"notfound"},
			expected:             inputText,
			found:                false,
			positionShouldChange: false,
		},
		{
			name:                 "Empty input",
			input:                []string{},
			expected:             inputText,
			found:                false,
			positionShouldChange: false,
		},
		{
			name:                 "Special symbols input",
			input:                []string{"test."},
			expected:             "hello world! This is a ",
			found:                true,
			positionShouldChange: true,
		},
		{
			name:           "Respect line breaks",
			input:          []string{"multistring"},
			expected:       `hello world! This is a test.`,
			found:          false,
			respectNewLine: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewLexer(inputText)

			previousPosition := tokenizer.pos

			result, found := tokenizer.consumeTill(test.respectNewLine, test.input...)

			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.found, found)
			assert.Equal(t, test.positionShouldChange, tokenizer.pos != previousPosition)
		})
	}
}
