package cmd

import (
	"github.com/specgen-io/specgen/specgen/v2/fail"
	"github.com/specgen-io/specgen/specgen/v2/markdown"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"strings"
)

func init() {
	cmdMarkdown.Flags().String(OutFile, "", OutFileDescription)
	cmdMarkdown.MarkFlagRequired(OutFile)

	rootCmd.AddCommand(cmdMarkdown)
}

var cmdMarkdown = &cobra.Command{
	Use:   "markdown",
	Short: "Generate markdown documentation",
	Run: func(cmd *cobra.Command, args []string) {
		outFile, err := cmd.Flags().GetString(OutFile)
		fail.IfError(err)

		err = GenMarkdownDocumentation(rootCmd, outFile, []string{"completion", "markdown"})
		fail.IfError(err, "Failed to write markdown documentation")
	},
}

func GenMarkdownDocumentation(cmd *cobra.Command, filePath string, skipCommands []string) error {
	m := markdown.NewMarkdown()
	GenRootCommandMarkdown(m, cmd, skipCommands)
	GenCommandsMarkdown(m, make([]string, 0), cmd, skipCommands)
	return ioutil.WriteFile(filePath, []byte(m.String()), 0644)
}

func GenRootCommandMarkdown(m *markdown.Markdown, command *cobra.Command, skipCommands []string) {
	command.InitDefaultHelpCmd()
	command.InitDefaultHelpFlag()

	name := command.CommandPath()

	m.Header1(name)
	m.Paragraph(command.Short)

	if len(command.Commands()) > 0 {
		m.Header2("Commands List")
		for _, subcommand := range command.Commands() {
			if !subcommand.IsAvailableCommand() || subcommand.IsAdditionalHelpTopicCommand() {
				continue
			}
			if contains(skipCommands, subcommand.Name()) {
				continue
			}
			m.ListItem(markdown.LinkHeader(subcommand.Name(), subcommand.Name()))
		}
		m.EmptyLine()
	}
}

func GenCommandsMarkdown(m *markdown.Markdown, commandPath []string, command *cobra.Command, skipCommands []string) {
	for _, subcommand := range command.Commands() {
		subcommandPath := append(commandPath, subcommand.Name())
		if !subcommand.IsAvailableCommand() || subcommand.IsAdditionalHelpTopicCommand() {
			continue
		}
		if contains(skipCommands, strings.Join(subcommandPath, " ")) {
			continue
		}
		GenCommandMarkdown(m, subcommandPath, subcommand)
		GenCommandsMarkdown(m, subcommandPath, subcommand, skipCommands)
	}
}

func GenCommandMarkdown(m *markdown.Markdown, cmdPath []string, command *cobra.Command) {
	command.InitDefaultHelpCmd()
	command.InitDefaultHelpFlag()

	name := strings.Join(cmdPath, " ")

	m.Header2(name)
	m.Paragraph(command.Short)

	if len(command.Long) > 0 {
		m.Paragraph(markdown.Bold(`Synopsis`))
		m.Paragraph(command.Long)
	}

	if command.Runnable() {
		m.CodeBlock(command.UseLine())
	}

	if len(command.Example) > 0 {
		m.Paragraph(markdown.Bold(`Examples`))
		m.CodeBlock(command.Example)
	}

	directFlags := command.NonInheritedFlags()
	parentFlags := command.InheritedFlags()

	if directFlags.HasAvailableFlags() || parentFlags.HasAvailableFlags() {
		m.Paragraph(markdown.Bold(`Flags`))

		flags := []*pflag.Flag{}
		if directFlags.HasAvailableFlags() {
			directFlags.VisitAll(func(flag *pflag.Flag) {
				flags = append(flags, flag)
			})
		}

		if parentFlags.HasAvailableFlags() {
			parentFlags.VisitAll(func(flag *pflag.Flag) {
				flags = append(flags, flag)
			})
		}

		for _, flag := range flags {
			if flagWithValue(flag) {
				writeFlag(m, flag)
			}
		}
		for _, flag := range flags {
			if !flagWithValue(flag) {
				writeFlag(m, flag)
			}
		}
	}
	m.EmptyLine()
}

func writeFlag(m *markdown.Markdown, flag *pflag.Flag) {
	command := markdown.Code("--" + flag.Name)
	if flag.Shorthand != "" {
		command = command + ", " + markdown.Code("-"+flag.Shorthand)
	}
	if flag.Usage != "" {
		command = command + " - " + firstLetterUpper(flag.Usage) + "."
	}
	m.ListItem(command)
	config := "Not required."
	if flagRequired(flag) {
		config = "Required."
	}
	if flag.DefValue != "" {
		config = config + " Default value: " + markdown.Code(flag.DefValue) + "."
	}
	m.Line("     " + config)
}

func firstLetterUpper(value string) string {
	return strings.ToUpper(value[0:1]) + value[1:]
}

func flagWithValue(flag *pflag.Flag) bool {
	return flag.Value.Type() != "bool"
}

func flagRequired(flag *pflag.Flag) bool {
	requiredAnnotation, found := flag.Annotations[cobra.BashCompOneRequiredFlag]
	return found && requiredAnnotation[0] == "true"
}
