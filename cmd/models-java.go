package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsJava.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsJava.Flags().String(PackageName, "", PackageNameDescription)
	cmdModelsJava.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsJava.MarkFlagRequired(SpecFile)
	cmdModelsJava.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsJava)
}

var cmdModelsJava = &cobra.Command{
	Use:   "models-java",
	Short: "Generate Java models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		packageName, err := cmd.Flags().GetString(PackageName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genjava.GenerateModels(specification, packageName, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write models code")
	},
}
