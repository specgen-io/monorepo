package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genopenapi"
)

func init() {
	cmdOpenapi.Flags().String(SpecFile, "", SpecFileDescription)
	cmdOpenapi.Flags().String(OutFile, "", OutFileDescription)

	cmdOpenapi.MarkFlagRequired(SpecFile)
	cmdOpenapi.MarkFlagRequired(OutFile)

	rootCmd.AddCommand(cmdOpenapi)
}

var cmdOpenapi = &cobra.Command{
	Use:   "openapi",
	Short: "Generate OpenAPI specification",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		outFile, err := cmd.Flags().GetString(OutFile)
		fail.IfError(err)

		err = genopenapi.GenerateSpecification(specFile, outFile)
		fail.IfErrorF(err, "Failed to generate OpeanAPI specifiction")
	},
}
