package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceScalaPlay.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceScalaPlay.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceScalaPlay.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceScalaPlay.Flags().String(ServicesPath, "", ServicesPathDescription)

	cmdServiceScalaPlay.MarkFlagRequired(SpecFile)
	cmdServiceScalaPlay.MarkFlagRequired(GeneratePath)

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

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		servicesPath, err := cmd.Flags().GetString(ServicesPath)
		fail.IfError(err)

		err = genscala.GeneratePlayService(specFile, swaggerPath, generatePath, servicesPath)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
