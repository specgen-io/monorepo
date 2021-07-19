package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceTsExpress.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceTsExpress.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceTsExpress.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceTsExpress.Flags().String(ServicesPath, "", ServicesPathDescription)
	cmdServiceTsExpress.Flags().String(Validation, "superstruct", ValidationDescription)

	cmdServiceTsExpress.MarkFlagRequired(SpecFile)

	cmdServiceTsExpress.MarkFlagRequired(SwaggerPath)
	cmdServiceTsExpress.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceTsExpress)

}

var cmdServiceTsExpress = &cobra.Command{
	Use:   "service-ts-express",
	Short: "Generate TypeScript Axios client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		swaggerPath, err := cmd.Flags().GetString(SwaggerPath)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		servicesPath, err := cmd.Flags().GetString(ServicesPath)
		fail.IfError(err)

		validation, err := cmd.Flags().GetString(Validation)
		fail.IfError(err)

		err = gents.GenerateExpressService(specFile, swaggerPath, generatePath, servicesPath, validation)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
