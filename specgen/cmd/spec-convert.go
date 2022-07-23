package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/convert/openapi"
	"github.com/specgen-io/specgen/v2/fail"
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
		fail.IfError(err)

		specFormat, err := cmd.Flags().GetString(Format)
		fail.IfError(err)

		outFile, err := cmd.Flags().GetString(OutFile)
		fail.IfError(err)

		if specFormat == "openapi" {
			err = openapi.ConvertFromOpenapi(inFile, outFile)
			fail.IfError(err)
		}

		console.PrintLnF(`Spec %s was successfully converted`, outFile)
	},
}
