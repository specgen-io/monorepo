package cmd

import (
	"bytes"
	"fmt"
	"github.com/specgen-io/specgen/v2/fail"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func init() {
	markdown.Flags().String(OutFile, "", OutFileDescription)
	markdown.MarkFlagRequired(OutFile)

	rootCmd.AddCommand(markdown)
}

var markdown = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	Run: func(cmd *cobra.Command, args []string) {
		outFile, err := cmd.Flags().GetString(OutFile)
		fail.IfError(err)

		GenMarkdownDocumentation(rootCmd, outFile)
	},
}

func GenCommandMarkdown(cmdPath []string, cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := strings.Join(cmdPath, " ")

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("**Synopsis**\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("**Examples**\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	flags := cmd.NonInheritedFlags()
	parentFlags := cmd.InheritedFlags()

	if flags.HasAvailableFlags() || parentFlags.HasAvailableFlags() {
		buf.WriteString("**Options**\n\n")
		buf.WriteString("```\n")
		if flags.HasAvailableFlags() {
			flags.SetOutput(buf)
			flags.PrintDefaults()
		}

		if parentFlags.HasAvailableFlags() {
			parentFlags.SetOutput(buf)
			parentFlags.PrintDefaults()
		}
		buf.WriteString("```\n\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func GenCommandsMarkdown(cmdPath []string, cmd *cobra.Command, w io.Writer) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if c.Name() == "completion" || c.Name() == "markdown" { //TODO: This is hardcode to omit certain commands
			continue
		}
		newCmdPath := append(cmdPath, c.Name())
		if err := GenCommandMarkdown(newCmdPath, c, w); err != nil {
			return err
		}
		if err := GenCommandsMarkdown(newCmdPath, c, w); err != nil {
			return err
		}
	}
	return nil
}

func GenRootCommandMarkdown(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("# " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")

	if len(cmd.Commands()) > 0 {
		buf.WriteString("## Commands List\n")
		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
				continue
			}
			if c.Name() == "completion" || c.Name() == "markdown" { //TODO: This is hardcode to omit certain commands
				continue
			}
			buf.WriteString("* [" + c.Name() + "](#" + c.Name() + ")\n")
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func GenMarkdownDocumentation(cmd *cobra.Command, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := GenRootCommandMarkdown(cmd, f); err != nil {
		return err
	}
	if err := GenCommandsMarkdown(make([]string, 0), cmd, f); err != nil {
		return err
	}
	return nil
}
