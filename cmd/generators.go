package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/spf13/cobra"
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
			params := map[string]string{}
			for _, arg := range generator.Args {
				value, err := cmd.Flags().GetString(arg.Name)
				fail.IfError(err)
				params[arg.Name] = value
			}
			specification := readSpecFile(params[gen.ArgSpecFile.Name])
			sources := generator.Generator(specification, params)
			err := sources.Write(false)
			fail.IfErrorF(err, "Failed to write source code")
		},
	}
	for _, arg := range generator.Args {
		command.Flags().String(arg.Name, arg.Default, arg.Description)
		if arg.Required {
			command.MarkFlagRequired(arg.Name)
		}
	}
	return command
}
