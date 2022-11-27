package generator

import (
	"generator/console"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"os"
	"sort"
	"spec"
	"strings"
)

func AddCobraCommands(parent *cobra.Command, generators []Generator) {
	for index := range generators {
		parent.AddCommand(generatorCommand(&generators[index]))
	}
}

func generatorCommand(g *Generator) *cobra.Command {
	command := &cobra.Command{
		Use:   g.Name,
		Short: g.Usage,
		Run: func(cmd *cobra.Command, args []string) {
			params := GeneratorArgsValues{}
			for _, arg := range g.Args {
				value, err := cmd.Flags().GetString(arg.Name)
				if err != nil {
					console.ProblemLn(err)
					os.Exit(1)
				}
				if arg.Values != nil {
					if !slices.Contains(arg.Values, value) {
						console.ProblemLnF(`Argument %s provided value "%s" is not among allowed: %s`, arg.Name, value, strings.Join(arg.Values, ", "))
						os.Exit(1)
					}
				}
				params[arg.Arg] = value
			}
			specification := readSpecFile(params[ArgSpecFile])
			sources := g.Generator(specification, params)
			err := sources.Write(false, func(wrote bool, fullpath string) {
				if wrote {
					console.PrintLn("Writing:", fullpath)
				} else {
					console.PrintLn("Skipping:", fullpath)
				}
			})
			if err != nil {
				console.ProblemLn("Failed to write source code")
				console.ProblemLn(err)
				os.Exit(1)
			}
		},
	}
	for _, arg := range g.Args {
		description := arg.Description
		if arg.Values != nil {
			description += "; allowed values: " + strings.Join(arg.Values, ", ")
		}
		command.Flags().String(arg.Name, arg.Default, description)
		if arg.Required {
			command.MarkFlagRequired(arg.Name)
		}
	}
	return command
}

func readSpecFile(specFile string) *spec.Spec {
	console.PrintLnF("Reading spec file: %s", specFile)
	data, err := ioutil.ReadFile(specFile)
	if err != nil {
		console.ProblemLnF("Failed to read spec file: %s", specFile)
		console.ProblemLn(err)
		os.Exit(1)
	}

	console.PrintLn("Parsing spec")
	specification, messages, err := spec.ReadSpec(data)

	if messages != nil {
		sort.Sort(messages.Items)

		for _, message := range messages.Items {
			if message.Level != spec.LevelError {
				console.PrintLnF(message.String())
			} else {
				console.ProblemLnF(message.String())
			}
		}
	}

	if err != nil {
		console.ProblemLnF("Failed to parse spec: %s", specFile)
		console.ProblemLn(err)
		os.Exit(1)
	}
	return specification
}
