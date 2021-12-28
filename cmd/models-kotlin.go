package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsKotlin.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsKotlin.Flags().String(PackageName, "", PackageNameDescription)
	cmdModelsKotlin.Flags().String(KotlinJsonLib, "jackson", KotlinJsonLibDescription)
	cmdModelsKotlin.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsKotlin.MarkFlagRequired(SpecFile)
	cmdModelsKotlin.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsKotlin)
}

var cmdModelsKotlin = &cobra.Command{
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

		sources := genkotlin.GenerateModels(specification, packageName, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write models code")
	},
}
