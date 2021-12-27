package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientGo.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientGo.Flags().String(ModuleName, "", ModuleNameDescription)
	cmdClientGo.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientGo.MarkFlagRequired(SpecFile)
	cmdClientGo.MarkFlagRequired(ModuleName)
	cmdClientGo.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientGo)
}

var cmdClientGo = &cobra.Command{
	Use:   "client-go",
	Short: "Generate Go client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		moduleName, err := cmd.Flags().GetString(ModuleName)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := gengo.GenerateGoClient(specification, moduleName, generatePath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write client source code")
	},
}
