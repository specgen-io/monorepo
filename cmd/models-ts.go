package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsTs.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsTs.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdModelsTs.Flags().String(Validation, "superstruct", ValidationDescription)

	cmdModelsTs.MarkFlagRequired(SpecFile)
	cmdModelsTs.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsTs)
}

var cmdModelsTs = &cobra.Command{
	Use:   "models-ts",
	Short: "Generate TypeScript models",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		validation, err := cmd.Flags().GetString(Validation)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		err = gents.GenerateModels(specification, generatePath, validation)
		fail.IfErrorF(err, "Failed to generate code")
	},
}
