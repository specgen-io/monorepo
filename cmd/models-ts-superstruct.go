package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/gents"
)

func init() {
	cmdModelsTsSuperstruct.Flags().String(SpecFile, "", SpecFileDescription)
	cmdModelsTsSuperstruct.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdModelsTsSuperstruct.MarkFlagRequired(SpecFile)
	cmdModelsTsSuperstruct.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdModelsTsSuperstruct)
}

var cmdModelsTsSuperstruct = &cobra.Command{
	Use:   "models-ts-superstruct",
	Short: "Generate TypeScript Superstruct models",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = gents.GenerateSuperstructModels(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate code")
	},
}