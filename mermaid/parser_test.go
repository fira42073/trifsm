package mermaid

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name       string
		inputState string
		result     *ASTNode
		wantErr    bool
	}{
		{
			name: "Simple",
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
			result: &ASTNode{
				Type:  "diagram",
				Value: DiagramHeader{},
				Children: []*ASTNode{
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "Still"},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "Crash",
							Description: "State of crash",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "Moving",
							Description: "Description of moving state",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name: KeywordTerminalState.String(),
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
							End: StateDeclValue{Name: "Still"},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "Still"},
							End: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "Still"},
							End: StateDeclValue{
								Name:        "Moving",
								Description: "Description of moving state",
							},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "Moving",
								Description: "Description of moving state",
							},
							End: StateDeclValue{Name: "Still"},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "Moving",
								Description: "Description of moving state",
							},
							End: StateDeclValue{
								Name:        "Crash",
								Description: "State of crash",
							},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "Crash",
								Description: "State of crash",
							},
							End: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
						},
					},
					{
						Type: ASTNodeTypeNote,
						Value: NoteValue{
							Placement: KeywordRight.String(),
							Ident:     "Still",
							Contents:  "This is a note",
						},
					},
				},
			},
		},
		{
			name: "More features",
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

	state Choice1 <<choice>>
    %% comment here for StateA
    StateA --> Choice1: Decision Point

	state ForkState <<fork>>
    StateB --> ForkState
    ForkState --> StateD
    ForkState --> StateE

	state JoinState <<join>>
    StateD --> JoinState
    StateE --> JoinState
    JoinState --> StateF

    StateF --> Choice1
    Choice1 --> StateG: Option 1
    Choice1 --> StateH: Option 2
    FinalState --> [*]

    state "Final State" as FinalState`,
			result: &ASTNode{
				Type:  "diagram",
				Value: DiagramHeader{},
				Children: []*ASTNode{
					{
						Type:  ASTNodeTypeDirection,
						Value: DirectionValue{Direction: "LR"},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateB"},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "StateC",
							Description: "C declr",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name: KeywordTerminalState.String(),
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "StateA",
							Description: "State A",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
							End: StateDeclValue{Name: "StateA"},
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "StateA",
							Description: "State A",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "StateA",
								Description: "State A",
							},
							End:         StateDeclValue{Name: "StateB"},
							Description: "Event1",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "StateB"},
							End: StateDeclValue{
								Name:        "StateC",
								Description: "C declr",
							},
							Description: "Event2",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "StateC",
								Description: "C declr",
							},
							End: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
						},
					},
					{
						Type: ASTNodeTypeNote,
						Value: NoteValue{
							Placement: KeywordRight.String(),
							Ident:     "StateA",
							Contents:  "This is a note right of StateA",
						},
					},
					{
						Type: ASTNodeTypeNote,
						Value: NoteValue{
							Placement: KeywordLeft.String(),
							Ident:     "StateB",
							Contents:  "This is a note left of StateB",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "Choice1",
							Description: KeywordChoice.String(),
						},
					},
					{
						Type: ASTNodeTypeComment,
						Value: CommentValue{
							Contents: "comment here for StateA",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "StateA",
								Description: "State A",
							},
							End: StateDeclValue{
								Name:        "Choice1",
								Description: KeywordChoice.String(),
							},
							Description: "Decision Point",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "ForkState",
							Description: KeywordFork.String(),
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "StateB"},
							End: StateDeclValue{
								Name:        "ForkState",
								Description: KeywordFork.String(),
							},
						},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateD"},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "ForkState",
								Description: KeywordFork.String(),
							},
							End: StateDeclValue{Name: "StateD"},
						},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateE"},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "ForkState",
								Description: KeywordFork.String(),
							},
							End: StateDeclValue{Name: "StateE"},
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "JoinState",
							Description: KeywordJoin.String(),
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "StateD"},
							End: StateDeclValue{
								Name:        "JoinState",
								Description: KeywordJoin.String(),
							},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "StateE"},
							End: StateDeclValue{
								Name:        "JoinState",
								Description: KeywordJoin.String(),
							},
						},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateF"},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "JoinState",
								Description: KeywordJoin.String(),
							},
							End: StateDeclValue{Name: "StateF"},
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "StateF"},
							End: StateDeclValue{
								Name:        "Choice1",
								Description: KeywordChoice.String(),
							},
						},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateG"},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "Choice1",
								Description: KeywordChoice.String(),
							},
							End:         StateDeclValue{Name: "StateG"},
							Description: "Option 1",
						},
					},
					{
						Type:  ASTNodeTypeStatedecl,
						Value: StateDeclValue{Name: "StateH"},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{
								Name:        "Choice1",
								Description: KeywordChoice.String(),
							},
							End:         StateDeclValue{Name: "StateH"},
							Description: "Option 2",
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "FinalState",
							Description: "Final State",
						},
					},
					{
						Type: ASTNodeTypeTransition,
						Value: TransitionValue{
							Start: StateDeclValue{Name: "FinalState"},
							End: StateDeclValue{
								Name: KeywordTerminalState.String(),
							},
						},
					},
					{
						Type: ASTNodeTypeStatedecl,
						Value: StateDeclValue{
							Name:        "FinalState",
							Description: "Final State",
						},
					},
				},
			},
		},

		{
			name: "Composite States",
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

	state Choice1 <<choice>>
    %% comment here for StateA
    StateA --> Choice1: Decision Point

	state ForkState <<fork>>
    StateB --> ForkState
    ForkState --> StateD
    ForkState --> StateE

	state JoinState <<join>>
    StateD --> JoinState
    StateE --> JoinState
    JoinState --> StateF

    StateF --> Choice1
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
			result:  nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create tokenizer and get tokens
			tokenizer := NewLexer(tc.inputState)
			tokens, err := tokenizer.Tokenize()
			require.NoError(t, err)

			diagram, err := NewASTBuilder(tokens).ParseDiagram()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !assert.Equal(t, tc.result, diagram) {
				t.Error(cmp.Diff(tc.result, diagram))
			}
		})
	}
}
