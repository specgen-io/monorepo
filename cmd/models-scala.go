package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsScala.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsScala.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsScala.MarkFlagRequired(SpecFile)
	cmdModelsScala.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsScala)
}

var cmdModelsScala = &cobra.Command{
	Use:   "models-scala",
	Short: "Generate Scala models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genscala.GenerateCirceModels(specification, "", generatePath)
		err = sources.Write(false)

		fail.IfErrorF(err, "Failed to generate models code")
	},
}
