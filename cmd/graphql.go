package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/gengql"
)

const OutputPath = "output-path"

func init() {
	cmdGraphql.Flags().String(SpecFile, "", SpecFileDescription)
	cmdGraphql.Flags().String(OutputPath, "", "Path to generated files")

	cmdGraphql.MarkFlagRequired(SpecFile)
	cmdGraphql.MarkFlagRequired(OutputPath)

	rootCmd.AddCommand(cmdGraphql)
}

var cmdGraphql = &cobra.Command{
	Use:   "graphql",
	Short: "Generate GraphQL schema",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		outputPath, err := cmd.Flags().GetString(OutputPath)
		fail.IfError(err)

		err = gengql.GenerateSchema(specFile, outputPath)
		fail.IfErrorF(err, "Failed to generate GraphQL schema")
	},
}
