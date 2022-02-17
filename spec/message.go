package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
)

type Level string

const (
	LevelError   Level = "error"
	LevelWarning Level = "warning"
	LevelInfo    Level = "info"
)

type Location struct {
	Line   int
	Column int
}

func locationFromNode(node *yaml.Node) *Location {
	if node == nil {
		return nil
	}
	return &Location{node.Line, node.Column}
}

type Message struct {
	Level    Level
	Message  string
	Location *Location
}

type Messages []Message

func (ms Messages) Len() int {
	return len(ms)
}
func (ms Messages) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
func (ms Messages) Less(i, j int) bool {
	if ms[i].Location == nil {
		return true
	}
	if ms[j].Location == nil {
		return false

	}
	return ms[i].Location.Line < ms[j].Location.Line ||
		(ms[i].Location.Line == ms[j].Location.Line && ms[i].Location.Column < ms[j].Location.Column)
}

func (ms Messages) Contains(level Level) bool {
	for _, m := range ms {
		if m.Level == level {
			return true
		}
	}
	return false
}
