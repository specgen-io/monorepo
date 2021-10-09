package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/spf13/cobra"
)

func init() {
	cmdServiceJavaSpring.Flags().String(SpecFile, "", SpecFileDescription)
	cmdServiceJavaSpring.Flags().String(PackageName, "", PackageNameDescription)
	cmdServiceJavaSpring.Flags().String(SwaggerPath, "", SwaggerPathDescription)
	cmdServiceJavaSpring.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdServiceJavaSpring.Flags().String(ServicesPath, "", ServicesPathDescription)

	cmdServiceJavaSpring.MarkFlagRequired(SpecFile)
	cmdServiceJavaSpring.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdServiceJavaSpring)
}

var cmdServiceJavaSpring = &cobra.Command{
	Use:   "service-java-spring",
	Short: "Generate Spring Java service source code",
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

		err = genjava.GenerateService(specification, packageName, swaggerPath, generatePath, servicesPath)
		fail.IfErrorF(err, "Failed to generate service code")
	},
}
