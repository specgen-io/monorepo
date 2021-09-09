package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	modelsJava.Flags().String(SpecFile, "", SpecFileDescription)
	modelsJava.Flags().String(GeneratePath, "", GeneratePathDescription)

	modelsJava.MarkFlagRequired(SpecFile)
	modelsJava.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(modelsJava)
}

var modelsJava = &cobra.Command{
	Use:   "models-java",
	Short: "Generate Java models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = genjava.GenerateModels(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
