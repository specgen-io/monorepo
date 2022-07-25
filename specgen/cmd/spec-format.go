package cmd

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/spf13/cobra"
	"io/ioutil"
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

		specification := readSpecFile(specFile)

		data, err := spec.WriteSpec(specification)
		fail.IfErrorF(err, "Failed to write spec")

		err = ioutil.WriteFile(specFile, data, 0644)
		fail.IfErrorF(err, "Failed to write spec to file: %s", specFile)
	},
}
