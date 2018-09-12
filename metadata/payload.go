package metadata

import (
	"strings"
)

const separator string = "<sep>"

type Tags map[string]interface{}

type Payload struct {
	useSeparator bool
	tags         Tags
}

func New() Payload {
	return Payload{tags: Tags{}}
}

func Separator() string {
	return separator
}

func (l *Payload) UseSeparator() bool {
	return l.useSeparator
}

func (l *Payload) Tags() Tags {
	return l.tags
}

func (l *Payload) AddBool(tag string, value bool) {
	l.tags[tag] = value
}

func (l *Payload) AddList(tag string, value []string) {
	var m = make(map[string]string)
	for _, v := range value {
		m[v] = v
	}

	var arr []string
	for v := range m {
		arr = append(arr, v)
	}

	l.tags[tag] = strings.Join(arr, separator)
	l.useSeparator = true
}

func (l *Payload) AddTag(tag string, value string) {
	l.tags[tag] = filter(value)
}

func filter(value interface{}) interface{} {
	if str, ok := value.(string); ok {
		if strings.EqualFold(str, "true") {
			return true
		}

		if strings.EqualFold(str, "false") {
			return false
		}

		return str
	}

	return value
}
