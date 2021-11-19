package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/spf13/cobra"
)

func init() {
	modelsKotlin.Flags().String(SpecFile, "", SpecFileDescription)
	modelsKotlin.Flags().String(PackageName, "", PackageNameDescription)
	modelsKotlin.Flags().String(GeneratePath, "", GeneratePathDescription)

	modelsKotlin.MarkFlagRequired(SpecFile)
	modelsKotlin.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(modelsKotlin)
}

var modelsKotlin = &cobra.Command{
	Use:   "models-kotlin",
	Short: "Generate Kotlin models source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		packageName, err := cmd.Flags().GetString(PackageName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		err = genkotlin.GenerateModels(specification, packageName, generatePath)
		fail.IfErrorF(err, "Failed to generate models code")
	},
}
