package cmd


import (
	"github.com/spf13/cobra"
	"specgen/fail"
	"specgen/gents"
)

func init() {
	cmdClientTsAxios.Flags().String(SpecFile, "", SpecFileDescription)
	cmdClientTsAxios.Flags().String(GeneratePath, "", GeneratePathDescription)

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

		err = gents.GenerateAxiosClient(specFile, generatePath)
		fail.IfErrorF(err, "Failed to generate code")
	},
}