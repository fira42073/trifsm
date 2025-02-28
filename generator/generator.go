package generator

import (
	"bytes"
	"fmt"
	"text/template"
)

type Generator struct {
	t *template.Template

	version      string
	inputCommand string
}

func NewGenerator(version, inputCommand string) *Generator {
	g := &Generator{
		t:            template.New("fsm.gotmpl"),
		version:      version,
		inputCommand: inputCommand,
	}
	g.loadTemplate()

	return g
}

func (g *Generator) Generate(fileName string) ([]byte, error) {
	f, err := g.parseFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("parsing file (%s): %w", fileName, err)
	}

	pkg := f.Name.Name

	m, err := g.extractMermaid(f)
	if err != nil {
		return nil, fmt.Errorf("extracting mermaid from comments: %w", err)
	}

	fsms, err := mermaidToFSMs(m)
	if err != nil {
		return nil, fmt.Errorf("parsing comments into mermaid: %w", err)
	}

	data := templateData{
		Command: g.inputCommand,
		Version: g.version,
		Package: pkg,
		FSMs:    fsms,
	}

	var b bytes.Buffer
	err = g.t.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	return b.Bytes(), nil
}
