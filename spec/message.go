package spec

import (
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"strconv"
	"strings"
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

func Error(messageFormat string, args ...interface{}) Message {
	return Message{LevelError, fmt.Sprintf(messageFormat, args...), nil}
}

func Warning(messageFormat string, args ...interface{}) Message {
	return Message{LevelWarning, fmt.Sprintf(messageFormat, args...), nil}
}

func Info(messageFormat string, args ...interface{}) Message {
	return Message{LevelInfo, fmt.Sprintf(messageFormat, args...), nil}
}

func (message Message) At(location *Location) Message {
	message.Location = location
	return message
}

func convertYamlError(err error, node *yaml.Node) Message {
	if strings.HasPrefix(err.Error(), "yaml: line ") {
		parts := strings.SplitN(strings.TrimPrefix(err.Error(), "yaml: line "), ":", 2)
		line, err := strconv.Atoi(parts[0])
		if err == nil {
			return Error(strings.TrimSpace(parts[1])).At(&Location{line, 0})
		}
	}
	return Error(err.Error()).At(locationFromNode(node))
}

type Messages struct {
	Items messages
}

func NewMessages() *Messages {
	return &Messages{[]Message{}}
}

func (ms *Messages) Add(message Message) {
	ms.Items = append(ms.Items, message)
}

func (ms *Messages) AddAll(messages ...Message) {
	ms.Items = append(ms.Items, messages...)
}

func (ms *Messages) ContainsLevel(level Level) bool {
	for _, m := range ms.Items {
		if m.Level == level {
			return true
		}
	}
	return false
}

func (ms *Messages) Contains(f func(Message) bool) bool {
	for _, m := range ms.Items {
		if f(m) {
			return true
		}
	}
	return false
}

type messages []Message

func (ms messages) Len() int {
	return len(ms)
}
func (ms messages) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
func (ms messages) Less(i, j int) bool {
	if ms[i].Location == nil {
		return true
	}
	if ms[j].Location == nil {
		return false

	}
	return ms[i].Location.Line < ms[j].Location.Line ||
		(ms[i].Location.Line == ms[j].Location.Line && ms[i].Location.Column < ms[j].Location.Column)
}
