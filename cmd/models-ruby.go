package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genruby"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsRuby.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsRuby.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsRuby.MarkFlagRequired(SpecFile)
	cmdModelsRuby.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsRuby)
}

var cmdModelsRuby = &cobra.Command{
	Use:   "models-ruby",
	Short: "Generate Ruby models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genruby.GenerateModels(specification, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write models code")
	},
}
