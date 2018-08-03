package metadata

import (
	"strconv"
	"strings"
)

const separator string = "<sep>"

type Payload struct {
	useSeparator bool
	tags         map[string]interface{}
}

func New() Payload {
	return Payload{tags: map[string]interface{}{}}
}

func Separator() string {
	return separator
}

func (l *Payload) UseSeparator() bool {
	return l.useSeparator
}

func (l *Payload) Tags() map[string]interface{} {
	return l.tags
}

func (l *Payload) AddTag(tag string, value interface{}) {
	l.tags[tag] = filter(value)

	if useSeparator(value) {
		l.useSeparator = true
	}
}

func filter(value interface{}) interface{} {
	arr, ok := value.([]string)
	if ok {
		return quote(strings.Join(arr, separator))
	}

	str, ok := value.(string)
	if ok {
		if strings.EqualFold(str, "true") {
			return true
		}

		if strings.EqualFold(str, "false") {
			return false
		}

		return quote(str)
	}

	return value
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func useSeparator(value interface{}) bool {
	_, ok := value.([]string)

	return ok
}
