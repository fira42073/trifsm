//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --lower --marshal --values
package mermaid

// ENUM(
// diagram, // Top-level diagram node
// direction, // Direction node
// statedecl, // State declaration
// transition, // Transition node
// note, // Note node
// composite_state, // Composite state node
// comment, // Comment node
// )
type ASTNodeType string

// DiagramValue is an interface for the values of the ASTNode
type DiagramValue interface{}

type Identifiable interface {
	Ident() string
}

type FilterableChildren []*ASTNode

func (f FilterableChildren) Filter(predicate func(*ASTNode) bool) []*ASTNode {
	nodes := make([]*ASTNode, 0, 10) // arbitrary initial capacity, should be enough for most usecases

	for _, node := range f {
		if predicate(node) {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (f FilterableChildren) FilterOne(predicate func(*ASTNode) bool) *ASTNode {
	for _, node := range f {
		if predicate(node) {
			return node
		}
	}

	return nil
}

// AST Node Definition
type ASTNode struct {
	Type     ASTNodeType
	Value    DiagramValue
	Children []*ASTNode
}

// FindEntrypointStateDecl either chooses the first node that is transition with start of terminalstate, or just the first state declaration. Always returns a state declaration node or nil.
func (n *ASTNode) FindEntrypointStateDecl() *ASTNode {
	// find first transition from terminal state
	startsWithTerminalState := FilterableChildren(n.Children).FilterOne(func(n *ASTNode) bool {
		if n.Type != ASTNodeTypeTransition {
			return false
		}

		transition, ok := n.Value.(TransitionValue)
		if !ok {
			return false
		}

		if transition.Start.Name == KeywordTerminalState.String() {
			return true
		}

		return false
	})

	// if we found a transition from terminal state, find the state declaration for that transition
	if startsWithTerminalState != nil {
		transition, ok := startsWithTerminalState.Value.(TransitionValue)
		if !ok {
			// This should never ever ever happen
			panic("transition value is not a TransitionValue")
		}

		return FilterableChildren(n.Children).FilterOne(func(a *ASTNode) bool {
			if n.Type != ASTNodeTypeStatedecl {
				return false
			}

			if decl, ok := a.Value.(StateDeclValue); ok && decl.Name == transition.End.Name {
				return true
			}

			return false
		})
	}

	// fuck it, just find the first state declaration
	return FilterableChildren(n.Children).FilterOne(func(n *ASTNode) bool {
		if n.Type != ASTNodeTypeStatedecl {
			return false
		}

		_, ok := n.Value.(StateDeclValue)
		return ok
	})
}

type DiagramHeader struct{}
type DirectionValue struct {
	Direction string
}

type TransitionValue struct {
	Start       StateDeclValue
	End         StateDeclValue
	Description string
}
type NoteValue struct {
	Placement string // KeywordLeft or KeywordRight
	Ident     string // of which state
	Contents  string
}

type CommentValue struct {
	Contents string
}
type StateDeclValue struct {
	Name        string
	Description string
}

func (s StateDeclValue) Ident() string {
	return s.Name
}

func NewDiagramNode() *ASTNode {
	return &ASTNode{
		Type:     ASTNodeTypeDiagram,
		Value:    DiagramHeader{},
		Children: nil,
	}
}
func NewDirectionNode(direction string) *ASTNode {
	return &ASTNode{
		Type: ASTNodeTypeDirection,
		Value: DirectionValue{
			Direction: direction,
		},
		Children: nil,
	}
}
func NewStateDeclNode(val StateDeclValue) *ASTNode {
	return &ASTNode{
		Type:     ASTNodeTypeStatedecl,
		Value:    val,
		Children: nil,
	}
}
func NewTransitionNode(val TransitionValue) *ASTNode {
	return &ASTNode{
		Type:     ASTNodeTypeTransition,
		Value:    val,
		Children: nil,
	}
}
func NewNoteNode(val NoteValue) *ASTNode {
	return &ASTNode{
		Type:  ASTNodeTypeNote,
		Value: val,
	}
}
func NewCommentNode(val string) *ASTNode {
	return &ASTNode{
		Type: ASTNodeTypeComment,
		Value: CommentValue{
			Contents: val,
		},
	}
}

var DiagramValueMap = map[ASTNodeType]DiagramValue{
	ASTNodeTypeDiagram:    DiagramHeader{},
	ASTNodeTypeDirection:  DirectionValue{},
	ASTNodeTypeStatedecl:  StateDeclValue{},
	ASTNodeTypeTransition: TransitionValue{},
	ASTNodeTypeNote:       NoteValue{},
	ASTNodeTypeComment:    CommentValue{},
}

var _ Identifiable = StateDeclValue{}
