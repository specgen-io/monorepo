package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientScalaHttps.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientScalaHttps.Flags().String(GeneratePath, "", GeneratePathDescription)

	cmdClientScalaHttps.MarkFlagRequired(SpecFile)
	cmdClientScalaHttps.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientScalaHttps)
}

var cmdClientScalaHttps = &cobra.Command{
	Use:   "client-scala-sttp",
	Short: "Generate Scala Sttp client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		err = genscala.GenerateSttpClient(specification, generatePath)
		fail.IfErrorF(err, "Failed to generate client code")
	},
}
