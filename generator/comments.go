package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// copyDocsToSpecs will take the GenDecl level documents and copy them
// to the children Type and Value specs.  I think this is actually working
// around a bug in the AST, but it works for now.
//
// Copied from https://github.com/abice/go-enum/blob/28240c662a5ec517eea3b58d6aa7743af31ae0a7/generator/generator.go#L730 , thanks <3
func copyGenDeclCommentsToSpecs(x *ast.GenDecl) {
	// Copy the doc spec to the type or value spec
	// cause they missed this... whoops
	if x.Doc != nil {
		for _, spec := range x.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			case *ast.ValueSpec:
				if s.Doc == nil {
					s.Doc = x.Doc
				}
			}
		}
	}
}

func containsFSM(ts *ast.TypeSpec) bool {
	if ts.Doc == nil {
		return false
	}

	for _, comment := range ts.Doc.List {
		if strings.Contains(comment.Text, `FSM`) {
			return true
		}
	}

	return false
}

// parseFile simply calls the parser.ParseFile function to parse comments.
func (g *Generator) parseFile(filename string) (*ast.File, error) {
	var fileSet = token.FileSet{}
	return parser.ParseFile(&fileSet, filename, nil, parser.ParseComments)
}

func (g *Generator) extractMermaid(f *ast.File) (map[string]string, error) {
	fsmComments := make(map[string]string)

	// inspired by https://github.com/abice/go-enum/blob/28240c662a5ec517eea3b58d6aa7743af31ae0a7/generator/generator.go#L698 , thanks <3
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			copyGenDeclCommentsToSpecs(node)
		case *ast.Ident:
			if node.Obj == nil || node.Obj.Kind != ast.Typ {
				return true // continue
			}

			// Make sure it's a spec (Type Identifiers can be throughout the code)
			// and that it contains the FSM keyword
			if ts, ok := node.Obj.Decl.(*ast.TypeSpec); ok && containsFSM(ts) {
				fsmComments[node.Name] = ts.Doc.Text()
			}
		}

		// Return true to continue through the tree
		return true
	})

	return fsmComments, nil
}
