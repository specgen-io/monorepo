package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/spf13/cobra"
)

var ClientGo = Generator{
	"client-go",
	"Generate Go client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *gen.Sources {
		return gengo.GenerateGoClient(specification, params[ModuleName], params[GeneratePath])
	},
}

var Generators = []Generator{
	ClientGo,
}

func createGeneratorCommand(generator *Generator) *cobra.Command {
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
			specification := readSpecFile(params[ArgSpecFile.Name])
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

func init() {
	for _, generator := range Generators {
		rootCmd.AddCommand(createGeneratorCommand(&generator))
	}
}
