package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientKotlin.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientKotlin.Flags().String(PackageName, "", PackageNameDescription)
	cmdClientKotlin.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientKotlin.MarkFlagRequired(SpecFile)
	cmdClientKotlin.MarkFlagRequired(PackageName)
	cmdClientKotlin.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientKotlin)
}

var cmdClientKotlin = &cobra.Command{
	Use:   "client-kotlin",
	Short: "Generate Kotlin client source code",
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
