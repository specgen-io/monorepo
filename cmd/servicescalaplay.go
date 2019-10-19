package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genscala"
)

func init() {
	cmdServiceScalaPlay.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceScalaPlay.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceScalaPlay.Flags().String(SourceManagedPath, "", SourceManagedPathDescription)
	cmdServiceScalaPlay.Flags().String(SourcePath, "", SourcePathDescription)
	cmdServiceScalaPlay.Flags().String(ResourcePath, "", ResourcePathDescription)

	cmdServiceScalaPlay.MarkFlagRequired(SpecFile)

	cmdServiceScalaPlay.MarkFlagRequired(SwaggerPath)
	cmdServiceScalaPlay.MarkFlagRequired(SourceManagedPath)
	cmdServiceScalaPlay.MarkFlagRequired(SourcePath)
	cmdServiceScalaPlay.MarkFlagRequired(ResourcePath)

	rootCmd.AddCommand(cmdServiceScalaPlay)
}

var cmdServiceScalaPlay = &cobra.Command{
	Use:   "service-scala-play",
	Short: "Generate Scala Play service source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		swaggerPath, err := cmd.Flags().GetString(SwaggerPath)
		fail.IfError(err)

		sourceManagedPath, err := cmd.Flags().GetString(SourceManagedPath)
		fail.IfError(err)

		sourcePath, err := cmd.Flags().GetString(SourcePath)
		fail.IfError(err)

		resourcePath, err := cmd.Flags().GetString(ResourcePath)
		fail.IfError(err)

		err = genscala.GeneratePlayService(specFile, swaggerPath, sourceManagedPath, sourcePath, resourcePath)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
