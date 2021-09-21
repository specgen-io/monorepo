package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientGo.Flags().String(ModuleName, "", ModuleNameDescription)

	cmdClientGo.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientGo.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientGo.MarkFlagRequired(SpecFile)
	cmdClientGo.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientGo)
}

var cmdClientGo = &cobra.Command{
	Use:   "client-go",
	Short: "Generate Go client source code",
	Run: func(cmd *cobra.Command, args []string) {
		moduleName, err := cmd.Flags().GetString(ModuleName)
		fail.IfError(err)

		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = gengo.GenerateGoClient(moduleName, specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate client code")
	},
}
