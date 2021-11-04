package cmd

import (
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/spec/old"
	"github.com/spf13/cobra"
)

func init() {
	cmdSpecConvert.Flags().String(SpecFile, "", SpecFileDescription)
	cmdSpecConvert.MarkFlagRequired(SpecFile)
	cmdSpecConvert.Flags().String(OutFile, "", OutFileDescription)
	rootCmd.AddCommand(cmdSpecConvert)
}

var cmdSpecConvert = &cobra.Command{
	Use:   "spec-convert",
	Short: "Convert spec from older versions to latest",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		outFile, err := cmd.Flags().GetString(OutFile)
		fail.IfError(err)

		if outFile == "" {
			outFile = specFile
		}

		err = old.FormatSpec(specFile, outFile, "2.1")
		fail.IfError(err)

		console.PrintLnF(`Spec %s was successfully converted`, outFile)
	},
}
