package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	for index := range gen.Generators {
		rootCmd.AddCommand(generatorCommand(&gen.Generators[index]))
	}
}

func generatorCommand(generator *gen.Generator) *cobra.Command {
	command := &cobra.Command{
		Use:   generator.Name,
		Short: generator.Usage,
		Run: func(cmd *cobra.Command, args []string) {
			params := map[gen.Arg]string{}
			for _, arg := range generator.Args {
				value, err := cmd.Flags().GetString(arg.Name)
				fail.IfError(err)
				if arg.Values != nil {
					if !contains(arg.Values, value) {
						fail.FailF(`Argument %s provided value "%s" is not among allowed: %s`, arg.Name, value, strings.Join(arg.Values, ", "))
					}
				}
				params[arg.Arg] = value
			}
			specification := readSpecFile(params[gen.ArgSpecFile])
			sources := generator.Generator(specification, params)
			err := sources.Write(false)
			fail.IfErrorF(err, "Failed to write source code")
		},
	}
	for _, arg := range generator.Args {
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
