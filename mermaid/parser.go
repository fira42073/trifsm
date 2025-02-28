package mermaid

import (
	"fmt"
	"strings"
)

// AST Builder
type ASTBuilder struct {
	tokens []Token
	pos    int
}

func NewASTBuilder(tokens []Token) *ASTBuilder {
	return &ASTBuilder{
		tokens: tokens,
		pos:    0,
	}
}

func (b *ASTBuilder) consume() Token {
	token := b.peek(0)
	b.pos++
	return token
}

func (b *ASTBuilder) consumeLine() []Token {
	tokens, lastPos := b.peekLine()
	b.pos = lastPos
	return tokens
}

func (b *ASTBuilder) peek(offset int) Token {
	return b.tokens[b.pos+offset]
}

func (b *ASTBuilder) peekLine() ([]Token, int) {
	startPos := b.pos
	defer func() {
		// Restore position
		b.pos = startPos
	}()

	tokens := make([]Token, 0)

	for {
		tok := b.consume() // temporarily advance to the new token
		if tok.Type == TokenTypeEOL || tok.Type == TokenTypeEOF {
			return tokens, b.pos
		}

		tokens = append(tokens, tok)
	}
}

func (b *ASTBuilder) isLastNode() bool {
	return b.pos >= len(b.tokens)
}

func (b *ASTBuilder) ParseDiagram() (*ASTNode, error) {
	if len(b.tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	root := NewDiagramNode()

	var inScope bool

	// Process the rest of the tokens
	for !b.isLastNode() {
		token := b.peek(0)

		switch token.Type {
		case TokenTypeKeyword:
			err := b.processKeyword(token, root, &inScope)
			if err != nil {
				return nil, err
			}

		case TokenTypeText:
			// %% this is come comment

			tokens := b.consumeLine()

			if len(tokens) != 1 && tokens[0].Type != TokenTypeText {
				return nil, fmt.Errorf("unexpected token number of tokens at position at line %d col %d", token.Line, token.Column)
			}

			root.Children = append(root.Children, NewCommentNode(strings.TrimSpace(tokens[0].Value)))

		case TokenTypeIdentifier:
			err := b.parseIdentifierStartingLine(root, b.consumeLine())
			if err != nil {
				return nil, fmt.Errorf("state transition: %w", err)
			}

		case TokenTypeEOL:
			// ignore dangling EOLs
			_ = b.consume()
			continue

		case TokenTypeEOF:
			return root, nil
		default:
			return nil, fmt.Errorf("unexpected token type: %v", token.Type)
		}
	}

	return root, nil
}

func (b *ASTBuilder) processKeyword(peekToken Token, root *ASTNode, inScope *bool) error {
	inScopeVal := *inScope
	defer func() {
		inScope = &inScopeVal
	}()

	switch Keyword(peekToken.Value) {
	case KeywordDirection:
		tokens := b.consumeLine() // Skip the direction keyword and the text after it
		if len(tokens) != 2 {
			return fmt.Errorf("unexpected number of tokens: %d", len(tokens))
		}

		if tokens[0].Type != TokenTypeKeyword && tokens[0].Value != KeywordDirection.String() {
			return fmt.Errorf("unexpected token type: %v", tokens[0].Type)
		}

		if tokens[1].Type != TokenTypeIdentifier {
			return fmt.Errorf("unexpected token type: %v", tokens[1].Type)
		}

		root.Children = append(root.Children, NewDirectionNode(tokens[1].Value))

	case KeywordState:
		tokens := b.consumeLine()

		node, err := b.parseStateKeywordStartingLine(peekToken, tokens)
		if err != nil {
			return fmt.Errorf("state keyword: %w", err)
		}

		// what we added just now is Identifiable
		// verify that there is no conflict and such node declaration doesn't already exist and redeclare if needed
		// last declaration wins
		if nodeDecl := nodeDeclarationExists(node, root); nodeDecl != nil {
			// redeclaration
			nodeDecl.Value = node.Value
		}

		root.Children = append(root.Children, node)

	case KeywordNote:
		// note right of Still: This is a note
		tokens := b.consumeLine()

		if len(tokens) != 6 {
			return fmt.Errorf("note declaration: unexpected number of tokens: %d", len(tokens))
		}

		if tokens[0].Value != KeywordNote.String() {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[0].Value)
		}

		if !(tokens[1].Value == KeywordRight.String() || tokens[1].Value == KeywordLeft.String()) {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[1].Value)
		}

		if tokens[2].Value != KeywordOf.String() {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[2].Value)
		}

		if tokens[3].Type != TokenTypeIdentifier {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[2].Value)
		}

		if tokens[4].Value != KeywordColon.String() {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[4].Value)
		}

		if tokens[5].Type != TokenTypeText {
			return fmt.Errorf("note declaration: unexpected token value on the line %d: %s", peekToken.Line, tokens[4].Value)
		}

		root.Children = append(root.Children, NewNoteNode(NoteValue{
			Placement: tokens[1].Value,
			Ident:     tokens[3].Value,
			Contents:  tokens[5].Value,
		}))

	case KeywordTerminalState:
		err := b.parseIdentifierStartingLine(root, b.consumeLine())
		if err != nil {
			return fmt.Errorf("state transition: %w", err)
		}

	case KeywordBlockStart:
		inScopeVal = true
	case KeywordBlockEnd:
		inScopeVal = false

	case KeywordConcurrent:

	default:
		return fmt.Errorf("unexpected keyword: %s", peekToken.Value)
	}
	return nil
}

func (b *ASTBuilder) parseStateKeywordStartingLine(token Token, tokens []Token) (*ASTNode, error) {
	var node *ASTNode

	switch len(tokens) {
	case 2:
		// state Still

		if tokens[1].Type != TokenTypeIdentifier {
			return nil, fmt.Errorf("unexpected token type: %v", tokens[3].Type)
		}

		node = NewStateDeclNode(StateDeclValue{
			Name:        tokens[1].Value,
			Description: "",
		})

	case 3:
		// state fork_state <<fork>>
		// state join_state <<join>>
		// state choice_state <<choice>>
		// state compositState {

		if tokens[2].Value == KeywordBlockStart.String() {
			return nil, fmt.Errorf("composite state declaration is not supported yet")
		}

		if tokens[1].Type != TokenTypeIdentifier {
			return nil, fmt.Errorf("unexpected token type: %v", tokens[1].Type)
		}

		// fork/join/choice keyword
		if tokens[2].Type != TokenTypeKeyword &&
			!(tokens[2].Value == KeywordFork.String() ||
				tokens[2].Value == KeywordJoin.String() ||
				tokens[2].Value == KeywordChoice.String()) {
			return nil, fmt.Errorf("unexpected token type: %v", tokens[2].Type)
		}

		switch tokens[2].Value {
		case KeywordFork.String():
			node = NewStateDeclNode(StateDeclValue{
				Name:        tokens[1].Value,
				Description: KeywordFork.String(),
			})
		case KeywordJoin.String():
			node = NewStateDeclNode(StateDeclValue{
				Name:        tokens[1].Value,
				Description: KeywordJoin.String(),
			})
		case KeywordChoice.String():
			node = NewStateDeclNode(StateDeclValue{
				Name:        tokens[1].Value,
				Description: KeywordChoice.String(),
			})
		}

	case 4:
		// state "This is a state description" as s2

		if tokens[1].Type != TokenTypeText {
			return nil, fmt.Errorf("unexpected token type: %v", tokens[1].Type)
		}

		if tokens[2].Value != KeywordAs.String() {
			return nil, fmt.Errorf("unexpected token value: %s", tokens[2].Value)
		}

		if tokens[3].Type != TokenTypeIdentifier {
			return nil, fmt.Errorf("unexpected token type: %v", tokens[3].Type)
		}

		node = NewStateDeclNode(StateDeclValue{
			Name:        tokens[3].Value,
			Description: tokens[1].Value,
		})

	default:
		return nil, fmt.Errorf("unexpected number of tokens on the line %d: %d", token.Line, len(tokens))
	}
	return node, nil
}

func nodeDeclarationExists(node *ASTNode, root *ASTNode) *ASTNode {
	newNode := node.Value.(Identifiable)
	// [*] is a terminal state, it doesn't need to be declared
	if newNode.Ident() == KeywordTerminalState.String() {
		return &ASTNode{
			Type: ASTNodeTypeStatedecl,
			Value: StateDeclValue{
				Name: KeywordTerminalState.String(),
			},
		}
	}

	existingState := FilterableChildren(root.Children).FilterOne(func(a *ASTNode) bool {
		identifiable, ok := a.Value.(Identifiable)
		if !ok {
			return false
		}

		return identifiable.Ident() == newNode.Ident()
	})

	return existingState
}

func (b *ASTBuilder) parseIdentifierStartingLine(root *ASTNode, tokens []Token) error {
	startIdentifier := tokens[0]
	if startIdentifier.Type != TokenTypeIdentifier {
		return fmt.Errorf("unexpected token type on line %d: %v", startIdentifier.Line, startIdentifier.Type)
	}

	switch len(tokens) {
	case 1:
		// SimpleStateDecl

		newNode := NewStateDeclNode(StateDeclValue{
			Name:        startIdentifier.Value,
			Description: "",
		})

		if nodeDecl := nodeDeclarationExists(newNode, root); nodeDecl != nil { // don't redeclare
			return nil
		}

		root.Children = append(root.Children, newNode)
	case 3:
		if tokens[1].Value == KeywordArrow.String() {
			// [*] --> IsPositive
			err := b.parseStateTransition(root, tokens)
			if err != nil {
				return err
			}
		} else if tokens[1].Value == KeywordColon.String() {
			// Still: this is a state
			newNode := NewStateDeclNode(StateDeclValue{
				Name:        startIdentifier.Value,
				Description: tokens[2].Value,
			})

			if nodeDecl := nodeDeclarationExists(newNode, root); nodeDecl != nil {
				// redeclaration
				nodeDecl.Value = StateDeclValue{
					Name:        startIdentifier.Value,
					Description: tokens[2].Value,
				}
			} else {
				// new one
				root.Children = append(root.Children, newNode)
			}

		} else {
			return fmt.Errorf("unexpected token value on line %d: %s", tokens[1].Line, tokens[1].Value)
		}

	case 5:
		// [*] --> IsPositive: gotta stay positive aight
		err := b.parseStateTransition(root, tokens)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *ASTBuilder) parseStateTransition(root *ASTNode, tokens []Token) error {
	// [*] --> IsPositive
	// IsNegative --> IsPositive

	startIdentifier := tokens[0]

	endIdentifier := tokens[2]
	if endIdentifier.Type != TokenTypeIdentifier {
		return fmt.Errorf("unexpected token type on line %d: %v", tokens[2].Line, tokens[2].Type)
	}

	// IsNegative --> IsPositive: gotta stay positive aight
	var stateTransitionDescription string
	if len(tokens) == 5 {
		if tokens[3].Value != KeywordColon.String() {
			return fmt.Errorf("unexpected token value on line %d: %s", tokens[3].Line, tokens[3].Value)
		}

		if tokens[4].Type != TokenTypeText {
			return fmt.Errorf("unexpected token type on line %d: %v", tokens[4].Line, tokens[4].Type)
		}

		stateTransitionDescription = tokens[4].Value
	}

	// find previously declared states and make sure that the state exists
	// if it doesn't then we need to add state declaration node as well
	startExistingState := FilterableChildren(root.Children).FilterOne(func(a *ASTNode) bool {
		identifiable, ok := a.Value.(Identifiable)
		if !ok {
			return false
		}

		return identifiable.Ident() == startIdentifier.Value
	})

	if startExistingState == nil {
		// create declaration first
		startExistingState = NewStateDeclNode(StateDeclValue{
			Name:        startIdentifier.Value,
			Description: "",
		})

		root.Children = append(root.Children, startExistingState)
	}

	// same for end state
	endExistingState := FilterableChildren(root.Children).FilterOne(func(a *ASTNode) bool {
		identifiable, ok := a.Value.(Identifiable)
		if !ok {
			return false
		}

		return identifiable.Ident() == endIdentifier.Value
	})
	if endExistingState == nil {
		// create declaration first
		endExistingState = NewStateDeclNode(StateDeclValue{
			Name:        endIdentifier.Value,
			Description: "",
		})

		root.Children = append(root.Children, endExistingState)

	}

	startIdent, err := digDiagramStateDecl(startExistingState.Value)
	if err != nil {
		return fmt.Errorf("digging diagram state declaration for start state: %w", err)
	}

	endIdent, err := digDiagramStateDecl(endExistingState.Value)
	if err != nil {
		return fmt.Errorf("digging diagram state declaration for end state: %w", err)
	}

	root.Children = append(root.Children, NewTransitionNode(TransitionValue{
		Start:       startIdent,
		End:         endIdent,
		Description: stateTransitionDescription,
	}))
	return nil
}

func digDiagramStateDecl(v DiagramValue) (StateDeclValue, error) {
	var endStateName, endStateDescription string
	switch typedState := v.(type) {
	case StateDeclValue:
		endStateName = typedState.Name
		endStateDescription = typedState.Description
	default:
		return StateDeclValue{}, fmt.Errorf("unidentifiable diagram value: %v", v)
	}

	return StateDeclValue{
		Name:        endStateName,
		Description: endStateDescription,
	}, nil
}
