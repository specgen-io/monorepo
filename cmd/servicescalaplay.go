package cmd

import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/genscala"
)

func init() {
	cmdServiceScalaPlay.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceScalaPlay.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceScalaPlay.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceScalaPlay.Flags().String(ServicesPath, "", ServicesPathDescription)
	cmdServiceScalaPlay.Flags().String(RoutesPath, "", RoutesPathDescription)

	cmdServiceScalaPlay.MarkFlagRequired(SpecFile)

	cmdServiceScalaPlay.MarkFlagRequired(SwaggerPath)
	cmdServiceScalaPlay.MarkFlagRequired(GeneratePath)
	cmdServiceScalaPlay.MarkFlagRequired(ServicesPath)
	cmdServiceScalaPlay.MarkFlagRequired(RoutesPath)

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

		routesPath, err := cmd.Flags().GetString(RoutesPath)
		fail.IfError(err)

		err = genscala.GeneratePlayService(specFile, swaggerPath, generatePath, servicesPath, routesPath)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
