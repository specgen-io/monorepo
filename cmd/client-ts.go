package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientTs.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientTs.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdClientTs.Flags().String(TsClient, "", TsClientDescription)
	cmdClientTs.Flags().String(TsValidation, "", TsValidationDescription)

	cmdClientTs.MarkFlagRequired(SpecFile)
	cmdClientTs.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientTs)
}

var cmdClientTs = &cobra.Command{
	Use:   "client-ts",
	Short: "Generate TypeScript client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		client, err := cmd.Flags().GetString(TsClient)
		fail.IfError(err)

		validation, err := cmd.Flags().GetString(TsValidation)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := gents.GenerateClient(specification, generatePath, client, validation)

		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write generate client code")
	},
}
