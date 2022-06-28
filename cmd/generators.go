package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/generators"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	for index := range generators.All {
		rootCmd.AddCommand(generatorCommand(&generators.All[index]))
	}
}

func generatorCommand(g *generator.Generator) *cobra.Command {
	command := &cobra.Command{
		Use:   g.Name,
		Short: g.Usage,
		Run: func(cmd *cobra.Command, args []string) {
			params := generator.GeneratorArgsValues{}
			for _, arg := range g.Args {
				value, err := cmd.Flags().GetString(arg.Name)
				fail.IfError(err)
				if arg.Values != nil {
					if !contains(arg.Values, value) {
						fail.FailF(`Argument %s provided value "%s" is not among allowed: %s`, arg.Name, value, strings.Join(arg.Values, ", "))
					}
				}
				params[arg.Arg] = value
			}
			specification := readSpecFile(params[generator.ArgSpecFile])
			sources := g.Generator(specification, params)
			err := sources.Write(false, func(wrote bool, fullpath string) {
				if wrote {
					console.PrintLn("Writing:", fullpath)
				} else {
					console.PrintLn("Skipping:", fullpath)
				}
			})
			fail.IfErrorF(err, "Failed to write source code")
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
