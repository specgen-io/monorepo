package cmd

import (
	"gopkg.in/specgen-io/specgen.v2/fail"
	"gopkg.in/specgen-io/specgen.v2/gengo"
	"github.com/spf13/cobra"
)

func init() {
	modelsGo.Flags().String(SpecFile, "", SpecFileDescription)
	modelsGo.Flags().String(GeneratePath, "", GeneratePathDescription)

	modelsGo.MarkFlagRequired(SpecFile)
	modelsGo.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(modelsGo)
}

var modelsGo = &cobra.Command{
	Use:   "models-go",
	Short: "Generate Go models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = gengo.GenerateModels(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
