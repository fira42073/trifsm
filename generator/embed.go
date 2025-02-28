package generator

import (
	"embed"
	"text/template"
)

//go:embed fsm.gotmpl
var content embed.FS

func (g *Generator) loadTemplate() {
	g.t = template.Must(g.t.ParseFS(content, "fsm.gotmpl"))
}
