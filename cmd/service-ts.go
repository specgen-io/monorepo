package cmd

import (
	"specgen/fail"
	"specgen/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceTs.Flags().String(SpecFile, "", SpecFileDescription)

	cmdServiceTs.Flags().String(TsServer, "", TsServerDescription)
	cmdServiceTs.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceTs.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceTs.Flags().String(ServicesPath, "", ServicesPathDescription)
	cmdServiceTs.Flags().String(Validation, "", ValidationDescription)

	cmdServiceTs.MarkFlagRequired(TsServer)
	cmdServiceTs.MarkFlagRequired(Validation)
	cmdServiceTs.MarkFlagRequired(SpecFile)
	cmdServiceTs.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceTs)

}

var cmdServiceTs = &cobra.Command{
	Use:   "service-ts",
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

		server, err := cmd.Flags().GetString(TsServer)
		fail.IfError(err)

		err = gents.GenerateService(specFile, swaggerPath, generatePath, servicesPath, server, validation)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
