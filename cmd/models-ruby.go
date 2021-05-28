package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genruby"
	"github.com/spf13/cobra"
)

func init() {
	modelsRuby.Flags().String(SpecFile, "", SpecFileDescription)
	modelsRuby.Flags().String(GeneratePath, "", GeneratePathDescription)

	modelsRuby.MarkFlagRequired(SpecFile)
	modelsRuby.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(modelsRuby)
}

var modelsRuby = &cobra.Command{
	Use:   "models-ruby",
	Short: "Generate Ruby models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = genruby.GenerateModels(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
