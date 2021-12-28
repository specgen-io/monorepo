package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientScala.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientScala.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientScala.MarkFlagRequired(SpecFile)
	cmdClientScala.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientScala)
}

var cmdClientScala = &cobra.Command{
	Use:   "client-scala",
	Short: "Generate Scala client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genscala.GenerateSttpClient(specification, "", generatePath)
		err = sources.Write(false)

		fail.IfErrorF(err, "Failed to generate client code")
	},
}
