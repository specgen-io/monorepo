package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientKotlinOkHttp.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientKotlinOkHttp.Flags().String(PackageName, "", PackageNameDescription)
	cmdClientKotlinOkHttp.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientKotlinOkHttp.MarkFlagRequired(SpecFile)
	cmdClientKotlinOkHttp.MarkFlagRequired(PackageName)
	cmdClientKotlinOkHttp.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientKotlinOkHttp)
}

var cmdClientKotlinOkHttp = &cobra.Command{
	Use:   "client-kotlin-okhttp",
	Short: "Generate OkHttp Kotlin client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		packageName, err := cmd.Flags().GetString(PackageName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genkotlin.GenerateClient(specification, packageName, generatePath)

		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write client code")
	},
}
