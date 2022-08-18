package cmd

import (
	"github.com/specgen-io/rendr/render"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const OutPath = "out"
const Set = "set"
const Values = "values"
const NoInput = "noinput"
const ForceInput = "forceinput"

func init() {
	cmdNew.Flags().String(OutPath, ".", `Path to output rendered template.`)
	cmdNew.Flags().StringArray(Set, []string{}, `Set arguments overrides in format "arg=value". Repeat for setting multiple arguments values.`)
	cmdNew.Flags().String(Values, "", `Path to arguments values JSON file.`)
	cmdNew.Flags().Bool(NoInput, false, `Do not request user input for missing arguments values.`)
	cmdNew.Flags().Bool(ForceInput, false, `Force user input requests even for noinput arguments.`)

	rootCmd.AddCommand(cmdNew)
}

var cmdNew = &cobra.Command{
	Use:   "new [template name]",
	Short: "Create new project with specgen code generation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]

		overrides, err := cmd.Flags().GetStringArray(Set)
		FailIfError(err)

		valuesJsonPath, err := cmd.Flags().GetString(Values)
		FailIfError(err)

		outPath, err := cmd.Flags().GetString(OutPath)
		FailIfError(err)

		noInput, err := cmd.Flags().GetBool(NoInput)
		FailIfError(err)

		forceInput, err := cmd.Flags().GetBool(ForceInput)
		FailIfError(err)

		inputMode := render.RegularInputMode
		if forceInput {
			inputMode = render.ForceInputMode
		}
		if noInput {
			inputMode = render.NoInputMode
		}

		var valuesJsonData []byte = nil
		if valuesJsonPath != "" {
			data, err := ioutil.ReadFile(valuesJsonPath)
			FailIfErrorF(err, `Failed to read arguments JSON file "%s"`, valuesJsonPath)
			valuesJsonData = data
		}

		err = renderTemplate(template, outPath, inputMode, valuesJsonData, overrides)
		FailIfError(err, "Failed to render template")
	},
}

func renderTemplate(name string, outPath string, inputMode render.InputMode, valuesJsonData []byte, overrides []string) error {
	template := render.Template{"https://github.com/specgen-io/templates", name, "rendr.yaml"}
	renderedFiles, err := template.Render(inputMode, valuesJsonData, overrides)
	if err != nil {
		return err
	}

	err = renderedFiles.WriteAll(outPath, true)
	if err != nil {
		return err
	}
	return nil
}
