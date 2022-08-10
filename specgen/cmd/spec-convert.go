package cmd

import (
	"generator/console"
	"specgen/convert/openapi"

	"github.com/spf13/cobra"
)

func init() {
	cmdSpecConvert.Flags().String(InFile, "", InFileDescription)
	cmdSpecConvert.MarkFlagRequired(InFile)
	cmdSpecConvert.Flags().String(Format, "", FormatDescription)
	cmdSpecConvert.MarkFlagRequired(Format)
	cmdSpecConvert.Flags().String(OutFile, "spec.yaml", OutFileDescription)
	rootCmd.AddCommand(cmdSpecConvert)
}

var cmdSpecConvert = &cobra.Command{
	Use:   "spec-convert",
	Short: "Convert spec from older versions to latest",
	Run: func(cmd *cobra.Command, args []string) {
		inFile, err := cmd.Flags().GetString(InFile)
		FailIfError(err)

		specFormat, err := cmd.Flags().GetString(Format)
		FailIfError(err)

		outFile, err := cmd.Flags().GetString(OutFile)
		FailIfError(err)

		if specFormat == "openapi" {
			err = openapi.ConvertFromOpenapi(inFile, outFile)
			FailIfError(err)
		}

		console.PrintLnF(`Spec %s was successfully converted`, outFile)
	},
}
