package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genruby"
)

func init() {
	cmdClientRuby.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientRuby.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientRuby.MarkFlagRequired(SpecFile)
	cmdClientRuby.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientRuby)
}

var cmdClientRuby = &cobra.Command{
	Use:   "client-ruby",
	Short: "Generate Ruby client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		err = genruby.GenerateClient(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate client code")
	},
}
