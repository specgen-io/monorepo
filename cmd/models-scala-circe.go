package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceScalaModels.Flags().String(SpecFile, "", SpecFileDescription)
	cmdServiceScalaModels.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdServiceScalaModels.MarkFlagRequired(SpecFile)
	cmdServiceScalaModels.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceScalaModels)
}

var cmdServiceScalaModels = &cobra.Command{
	Use:   "models-scala-circe",
	Short: "Generate Scala models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		err = genscala.GenerateServiceModels(specification, "", generatePath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
