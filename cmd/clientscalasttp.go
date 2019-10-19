package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genscala"
)

func init() {
	cmdClientScalaHttps.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientScalaHttps.Flags().String(SourceManagedPath, "", SourceManagedPathDescription)

	cmdClientScalaHttps.MarkFlagRequired(SpecFile)
	cmdClientScalaHttps.MarkFlagRequired(SourceManagedPath)

	rootCmd.AddCommand(cmdClientScalaHttps)
}

var cmdClientScalaHttps = &cobra.Command{
	Use:   "client-scala-sttp",
	Short: "Generate Scala Sttp client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		sourceManagedPath, err := cmd.Flags().GetString(SourceManagedPath)
		fail.IfError(err)

		err = genscala.GenerateSttpClient(specFile, sourceManagedPath)
		fail.IfErrorF(err, "Failed to generate client code")
	},
}
