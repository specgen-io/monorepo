package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceScala.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceScala.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceScala.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceScala.Flags().String(ServicesPath, "", ServicesPathDescription)

	cmdServiceScala.MarkFlagRequired(SpecFile)
	cmdServiceScala.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceScala)
}

var cmdServiceScala = &cobra.Command{
	Use:   "service-scala",
	Short: "Generate Scala service source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		swaggerPath, err := cmd.Flags().GetString(SwaggerPath)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		servicesPath, err := cmd.Flags().GetString(ServicesPath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genscala.GeneratePlayService(specification, swaggerPath, generatePath, servicesPath)
		err = sources.Write(false)

		fail.IfErrorF(err, "Failed to generate service code")
	},
}
