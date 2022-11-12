package markdown

import "fmt"

func LinkHeader(name, reference string) string {
	return fmt.Sprintf(`[%s](#%s)`, name, reference)
}

func Bold(value string) string {
	return "**" + value + "**"
}

func Code(value string) string {
	return "`" + value + "`"
}
