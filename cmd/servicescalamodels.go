package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genscala"
)

func init() {
	cmdServiceScalaModels.Flags().String(SpecFile, "", SpecFileDescription)
	cmdServiceScalaModels.Flags().String(SourceManagedPath, "", "Path to managed source code files")

	cmdServiceScalaModels.MarkFlagRequired(SpecFile)
	cmdServiceScalaModels.MarkFlagRequired(SourceManagedPath)

	rootCmd.AddCommand(cmdServiceScalaModels)
}

var cmdServiceScalaModels = &cobra.Command{
	Use:   "service-scala-models",
	Short: "Generate Scala models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		sourceManagedPath, err := cmd.Flags().GetString(SourceManagedPath)
		fail.IfError(err)

		err = genscala.GenerateServiceModels(specFile, sourceManagedPath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
