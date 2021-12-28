package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientJava.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientJava.Flags().String(PackageName, "", PackageNameDescription)
	cmdClientJava.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientJava.MarkFlagRequired(SpecFile)
	cmdClientJava.MarkFlagRequired(PackageName)
	cmdClientJava.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientJava)
}

var cmdClientJava = &cobra.Command{
	Use:   "client-java",
	Short: "Generate Java client source code",
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
