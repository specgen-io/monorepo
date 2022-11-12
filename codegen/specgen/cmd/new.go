package cmd

import (
	"fmt"
	"github.com/specgen-io/rendr/render"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const OutPath = "out"
const Set = "set"
const ExtraRoots = "root"
const Values = "values"
const NoInput = "noinput"
const ForceInput = "forceinput"
const Source = "source"
const NoOverwrites = "nooverwrites"

func init() {
	cmdNew.Flags().String(OutPath, ".", `path to output rendered template`)
	cmdNew.Flags().StringArray(Set, []string{}, `set arguments overrides in format "arg=value", repeat for setting multiple arguments values`)
	cmdNew.Flags().String(Values, "", `path to arguments values JSON file`)
	cmdNew.Flags().String(Source, "https://github.com/specgen-io/templates.git", `location of templates`)
	cmdNew.Flags().Bool(NoInput, false, `do not request user input for missing arguments values`)
	cmdNew.Flags().Bool(ForceInput, false, `force user input requests even for noinput arguments`)
	cmdNew.Flags().Bool(NoOverwrites, false, `do not overwrite files with rendered from template`)
	cmdNew.Flags().StringArray(ExtraRoots, []string{}, `extra template root, repeat for setting multiple extra roots`)

	rootCmd.AddCommand(cmdNew)
}

var cmdNew = &cobra.Command{
	Use:   "new [template name]",
	Short: "Create new project with specgen code generation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]

		outPath, err := cmd.Flags().GetString(OutPath)
		FailIfError(err)

		overrides, err := cmd.Flags().GetStringArray(Set)
		FailIfError(err)

		valuesJsonPath, err := cmd.Flags().GetString(Values)
		FailIfError(err)

		source, err := cmd.Flags().GetString(Source)
		FailIfError(err)

		noInput, err := cmd.Flags().GetBool(NoInput)
		FailIfError(err)

		forceInput, err := cmd.Flags().GetBool(ForceInput)
		FailIfError(err)

		noOverwrites, err := cmd.Flags().GetBool(NoOverwrites)
		FailIfError(err)

		extraRoots, err := cmd.Flags().GetStringArray(ExtraRoots)
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

		sourceUrl := fmt.Sprintf(`%s/%s`, source, template)
		err = renderTemplate(sourceUrl, extraRoots, outPath, inputMode, valuesJsonData, overrides, !noOverwrites)
		FailIfError(err, "Failed to render template")
	},
}

func renderTemplate(sourceUrl string, extraRoots []string, outPath string, inputMode render.InputMode, valuesJsonData []byte, overrides []string, overwriteFiles bool) error {
	template := render.Template{sourceUrl, "rendr.yaml", extraRoots}
	renderedFiles, err := template.Render(inputMode, valuesJsonData, overrides)
	if err != nil {
		return err
	}

	err = renderedFiles.WriteAll(outPath, overwriteFiles)
	if err != nil {
		return err
	}
	return nil
}
