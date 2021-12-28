package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceJava.Flags().String(SpecFile, "", SpecFileDescription)
	cmdServiceJava.Flags().String(PackageName, "", PackageNameDescription)
	cmdServiceJava.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceJava.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceJava.Flags().String(ServicesPath, "", ServicesPathDescription)

	cmdServiceJava.MarkFlagRequired(SpecFile)
	cmdServiceJava.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceJava)
}

var cmdServiceJava = &cobra.Command{
	Use:   "service-java",
	Short: "Generate Java service source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		packageName, err := cmd.Flags().GetString(PackageName)
		fail.IfError(err)

		swaggerPath, err := cmd.Flags().GetString(SwaggerPath)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		servicesPath, err := cmd.Flags().GetString(ServicesPath)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		sources := genjava.GenerateService(specification, packageName, swaggerPath, generatePath, servicesPath)
		err = sources.Write(false)
		fail.IfErrorF(err, "Failed to write service code")
	},
}
