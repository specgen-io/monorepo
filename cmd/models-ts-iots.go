package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdModelsTsIoTs.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsTsIoTs.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsTsIoTs.MarkFlagRequired(SpecFile)
	cmdModelsTsIoTs.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsTsIoTs)
}

var cmdModelsTsIoTs = &cobra.Command{
	Use:   "models-ts-iots",
	Short: "Generate TypeScript io-ts models",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = gents.GenerateIoTsModels(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate code")
	},
}