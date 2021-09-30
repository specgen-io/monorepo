package cmd

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/spf13/cobra"
)

func init() {
	cmdSpecFormat.Flags().String(SpecFile, "", SpecFileDescription)
	cmdSpecFormat.MarkFlagRequired(SpecFile)
	rootCmd.AddCommand(cmdSpecFormat)
}

var cmdSpecFormat = &cobra.Command{
	Use:   "spec-format",
	Short: "Format spec",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		result, err := spec.ReadSpecFile(specFile)
		fail.IfErrorF(err, "Failed to read spec")

		err = spec.WriteSpecFile(result.Spec, specFile)
		fail.IfErrorF(err, "Failed to write spec")
	},
}
