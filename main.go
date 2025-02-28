package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"
	"strings"

	"github.com/fira42073/trifsm/generator"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
)

var (
	version string
)

type rootT struct {
	OutputSuffix      string
	FileNames         cli.StringSlice
	TemplateFileNames cli.StringSlice
	Aliases           cli.StringSlice
	BuildTags         cli.StringSlice
}

func main() {
	var argv rootT

	clr := color.New()
	out := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(clr.Output(), format, args...)
	}

	app := &cli.App{
		Name:            "trifsm",
		Usage:           "mermaid diagram to fsm",
		HideHelpCommand: true,
		Version:         version,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				EnvVars:     []string{"GOFILE"},
				Usage:       "The file(s) to generate FSMs.  Use more than one flag for more files.",
				Required:    true,
				Destination: &argv.FileNames,
			},
			&cli.StringFlag{
				Name:        "output-suffix",
				Usage:       "Changes the default filename suffix of _fsm to something else.  `.go` will be appended to the end of the string no matter what, so that `_test.go` cases can be accommodated ",
				Destination: &argv.OutputSuffix,
			},
		},
		Action: func(ctx *cli.Context) error {
			g := generator.NewGenerator(version, "trifsm")

			for _, fileOption := range argv.FileNames.Value() {

				var filenames []string
				if fn, err := globFilenames(fileOption); err != nil {
					return err
				} else {
					filenames = fn
				}

				outputSuffix := `_fsm`
				if argv.OutputSuffix != "" {
					outputSuffix = argv.OutputSuffix
				}

				for _, fileName := range filenames {
					originalName := fileName

					out("trifsm processing: %s\n", color.Cyan(originalName))
					fileName, _ = filepath.Abs(fileName)

					outFilePath := fmt.Sprintf("%s%s.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)), outputSuffix)
					if strings.HasSuffix(fileName, "_test.go") {
						outFilePath = strings.Replace(outFilePath, "_test"+outputSuffix+".go", outputSuffix+"_test.go", 1)
					}

					// Parse the file given in arguments
					raw, err := g.Generate(fileName)
					if err != nil {
						return fmt.Errorf("failed generating FSMs\nInputFile=%s\nError=%s", color.Cyan(fileName), color.RedBg(err))
					}

					// Nothing was generated, ignore the output and don't create a file.
					if len(raw) < 1 {
						out(color.Yellow("trifsm ignored. file: %s\n"), color.Cyan(originalName))
						continue
					}

					mode := int(0o644)
					err = os.WriteFile(outFilePath, raw, os.FileMode(mode))
					if err != nil {
						return fmt.Errorf("failed writing to file %s: %s", color.Cyan(outFilePath), color.Red(err))
					}
					out("trifsm finished. file: %s\n", color.Cyan(originalName))
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// globFilenames gets a list of filenames matching the provided filename.
// In order to maintain existing capabilities, only glob when a * is in the path.
// Leave execution on par with old method in case there are bad patterns in use that somehow
// work without the Glob method.
func globFilenames(filename string) ([]string, error) {
	if strings.Contains(filename, "*") {
		matches, err := filepath.Glob(filename)
		if err != nil {
			return []string{}, fmt.Errorf("failed parsing glob filepath\nInputFile=%s\nError=%s", color.Cyan(filename), color.RedBg(err))
		}
		return matches, nil
	} else {
		return []string{filename}, nil
	}
}
