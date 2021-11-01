package cmd

import (
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/spf13/cobra"
)

func init() {
	cmdClientTsAxios.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientTsAxios.Flags().String(GeneratePath, "", GeneratePathDescription)
	cmdClientTsAxios.Flags().String(Validation, "", ValidationDescription)

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

		validation, err := cmd.Flags().GetString(Validation)
		fail.IfError(err)

		specification := readSpecFile(specFile)

		err = gents.GenerateAxiosClient(specification, generatePath, validation)
		fail.IfErrorF(err, "Failed to generate code")
	},
}
