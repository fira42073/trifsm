package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fira42073/trifsm/mermaid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// add type prefix and
func identVar(typ string, stateName string) string {
	return typ + cases.Title(language.English, cases.Compact).String(stateName)
}

func diagramToStateDecl(typeDecl string, root *mermaid.ASTNode) (initialStateName string, decls []stateDecl, err error) {
	usd := newUniqueStateDecls(len(root.Children))

	for _, child := range root.Children {
		if child.Type != mermaid.ASTNodeTypeStatedecl {
			continue // skip non-state declarations
		}

		stateDeclValue, ok := child.Value.(mermaid.StateDeclValue)
		if !ok {
			return "", nil, fmt.Errorf("unexpected value type: %T", child.Value)
		}

		usd.add(stateDecl{
			Name:    identVar(typeDecl, stateDeclValue.Ident()),
			Type:    typeDecl,
			Value:   stateDeclValue.Ident(),
			Comment: stateDeclValue.Description,
		})
	}

	entrypointStateNode := root.FindEntrypointStateDecl()
	if entrypointStateNode == nil {
		return "", nil, fmt.Errorf("no entrypoint state found")
	}

	return identVar(typeDecl, entrypointStateNode.Value.(mermaid.StateDeclValue).Ident()), usd.collect(), nil
}

func diagramToTransitionEvents(typeDecl string, root *mermaid.ASTNode) (transitionsTemplateValue string, transitionEvents []transitionEvent, err error) {
	transitionEvents = make([]transitionEvent, 0, len(root.Children))

	var events strings.Builder
	events.WriteString("\tfsm.Events{\n")

	for _, child := range root.Children {
		if child.Type != mermaid.ASTNodeTypeTransition {
			continue // skip non-transitions
		}

		stateTransitionValue, ok := child.Value.(mermaid.TransitionValue)
		if !ok {
			return "", nil, fmt.Errorf("unexpected value type: %T", child.Value)
		}

		startConst := identVar(typeDecl, stateTransitionValue.Start.Ident())
		endConst := identVar(typeDecl, stateTransitionValue.End.Ident())

		// Opening bracket
		events.WriteString("\t\t\t{")

		// Name
		events.WriteString("Name: ")

		teName := startConst + "__" + endConst
		teVal := startConst + `+"__"+` + endConst

		transitionEvents = append(transitionEvents, transitionEvent{
			Name:    teName,
			Type:    typeDecl,
			Value:   teVal,
			Comment: stateTransitionValue.Description,
		})

		events.WriteString(teName)
		events.WriteString(", ")

		// Src
		events.WriteString("Src: []string{")
		events.WriteString(startConst)
		events.WriteString("}, ")

		// Dst
		events.WriteString("Dst: ")
		events.WriteString(endConst)

		// Closing bracket
		events.WriteString("},\n")
	}

	events.WriteString("\t\t}")
	return events.String(), transitionEvents, nil
}

// normalize string of type
// `some comments here
// more comments about the type
//
// FSM
// ---
// stateDiagram-v2
// ....`
// to trim leading FSM --- and spaces
func trimMermaidComment(comment string) string {
	lines := strings.Split(comment, "\n")
	startIdx := -1

	// Find the first occurrence of "FSM" on its own line
	for i, line := range lines {
		if strings.TrimSpace(line) == "FSM" {
			startIdx = i
			break
		}
	}

	if startIdx == -1 || startIdx+1 >= len(lines) {
		panic("no FSM found")
	}

	// Trim the "---" line after "FSM" and return the rest
	trimmed := strings.Join(lines[startIdx+2:], "\n")
	return strings.TrimSpace(trimmed)
}

func fsmType(s string) string {
	return s + "FSM"
}

func mermaidToFSMs(mermaids map[string]string) ([]fsmData, error) {
	fsms := make([]fsmData, 0, len(mermaids))

	for typ, val := range mermaids {
		val = trimMermaidComment(val)
		lex := mermaid.NewLexer(val)

		tokens, err := lex.Tokenize()
		if err != nil {
			return nil, fmt.Errorf("tokenizing mermaid: %w", err)
		}

		rootDiagramNode, err := mermaid.NewASTBuilder(tokens).ParseDiagram()
		if err != nil {
			return nil, fmt.Errorf("parsing mermaid: %w", err)
		}

		initialStateName, states, err := diagramToStateDecl(fsmType(typ), rootDiagramNode)
		if err != nil {
			return nil, fmt.Errorf("parsing states: %w", err)
		}

		transitionTemplateValue, transitionEventDefs, err := diagramToTransitionEvents(fsmType(typ), rootDiagramNode)
		if err != nil {
			return nil, fmt.Errorf("parsing transitions: %w", err)
		}

		fsms = append(fsms, fsmData{
			CustomTypeName:      fsmType(typ),
			ConstructorName:     "New" + fsmType(typ),
			InitialState:        initialStateName,
			TransitionEvents:    transitionTemplateValue,
			TransitionEventDefs: transitionEventDefs,
			States:              states,
		})
	}

	// deterministic order
	sort.Sort(sortedFSMs(fsms))

	return fsms, nil
}
