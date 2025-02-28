//go:generate go run github.com/abice/go-enum@v0.6.0 -f=$GOFILE --lower --marshal --values

package mermaid

// TokenType represents different types of tokens in the state diagram
// ENUM(
// Keyword=0
// Text
// Identifier
// BlockStart
// BlockEnd
// EOF
// EOL
// )
type TokenType int

const (
	kw = Keyword("")
)

// ENUM(
// state
// direction
// note
// as
// right
// left
// of
// concurrent=--
// colon=:
// arrow=-->
// fork=<<fork>>
// join=<<join>>
// choice=<<choice>>
// comment=%%
// blockStart={
// blockEnd=}
// terminalState=[*]
// )
type Keyword string

func (kw Keyword) ToStringFunc() func(k Keyword) string {
	return func(k Keyword) string {
		return k.String()
	}
}

// ENUM(
// stateDiagramV1="stateDiagram"
// stateDiagramV2="stateDiagram-v2"
// )
type HeaderType string

type HeaderInfo struct {
	Type  HeaderType
	Title string
}

// Token represents a single token from the input
type Token struct {
	Value  string
	Type   TokenType
	Line   uint
	Column uint
}
