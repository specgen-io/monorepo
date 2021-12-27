package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientJavaOkHttp.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientJavaOkHttp.Flags().String(PackageName, "", PackageNameDescription)
	cmdClientJavaOkHttp.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientJavaOkHttp.MarkFlagRequired(SpecFile)
	cmdClientJavaOkHttp.MarkFlagRequired(PackageName)
	cmdClientJavaOkHttp.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientJavaOkHttp)
}

var cmdClientJavaOkHttp = &cobra.Command{
	Use:   "client-java-okhttp",
	Short: "Generate OkHttp Java client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		packageName, err := cmd.Flags().GetString(PackageName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genjava.GenerateClient(specification, packageName, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write client code")
	},
}
