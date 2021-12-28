package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsGo.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsGo.Flags().String(ModuleName, "", ModuleNameDescription)
	cmdModelsGo.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsGo.MarkFlagRequired(SpecFile)
	cmdModelsGo.MarkFlagRequired(ModuleName)
	cmdModelsGo.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsGo)
}

var cmdModelsGo = &cobra.Command{
	Use:   "models-go",
	Short: "Generate Go models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		moduleName, err := cmd.Flags().GetString(ModuleName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := gengo.GenerateModels(specification, moduleName, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write models code")
	},
}
