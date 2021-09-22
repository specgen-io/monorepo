package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceGo.Flags().String(SpecFile, "", SpecFileDescription)
	cmdServiceGo.Flags().String(ModuleName, "", ModuleNameDescription)
	cmdServiceGo.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceGo.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceGo.Flags().String(ServicesPath, "", ServicesPathDescription)

	cmdServiceGo.MarkFlagRequired(SpecFile)
	cmdServiceGo.MarkFlagRequired(ModuleName)
	cmdServiceGo.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceGo)
}

var cmdServiceGo = &cobra.Command{
	Use:   "service-go",
	Short: "Generate Go service source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		moduleName, err := cmd.Flags().GetString(ModuleName)
		fail.IfError(err)

		swaggerPath, err := cmd.Flags().GetString(SwaggerPath)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		servicesPath, err := cmd.Flags().GetString(ServicesPath)
		fail.IfError(err)

		err = gengo.GenerateService(specFile, moduleName, swaggerPath, generatePath, servicesPath)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
