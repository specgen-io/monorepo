package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientTsAxios.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientTsAxios.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdClientTsAxios.Flags().String("ts-models", "io-ts", "TypeScript models library")

	cmdClientTsAxios.MarkFlagRequired(SpecFile)
	cmdClientTsAxios.MarkFlagRequired(GeneratePath)

	rootCmd.AddCommand(cmdClientTsAxios)
}

var cmdClientTsAxios = &cobra.Command{
	Use:   "client-ts-axios",
	Short: "Generate TypeScript Axios client source code",
	Run: func(cmd *cobra.Command, args []string) {
		specFile, err := cmd.Flags().GetString(SpecFile)
		fail.IfError(err)

		generatePath, err := cmd.Flags().GetString(GeneratePath)
		fail.IfError(err)

		tsModels, err := cmd.Flags().GetString("ts-models")
		fail.IfError(err)

		err = gents.GenerateAxiosClient(specFile, generatePath, tsModels)
		fail.IfErrorF(err, "Failed to generate code")
	},
}
