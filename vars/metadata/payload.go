package metadata

import (
	"github.com/vvval/go-metadata-scanner/util"
	"strings"
)

const separator string = "<sep>"

type Tags map[string]interface{}

func (t Tags) Tag(tag string) (interface{}, bool) {
	v, ok := t[tag]

	return v, ok
}

func (t Tags) Count() int {
	return len(t)
}

type Payload struct {
	tags  Tags
	lists []string
}

func New() Payload {
	return Payload{tags: Tags{}}
}

func Separator() string {
	return separator
}

func (p *Payload) UseSeparator() bool {
	for _, tag := range p.lists {
		if strings.Contains(p.tags[tag].(string), separator) {
			return true
		}
	}

	return false
}

func (p *Payload) Tags() Tags {
	return p.tags
}

func (p *Payload) ListTags() []string {
	return p.lists
}

func (p *Payload) AddBool(tag string, value bool) {
	p.tags[tag] = value
}

func (p *Payload) AddList(tag string, value []string) {
	p.tags[tag] = strings.Join(util.UniqueValues(value), separator)
	p.lists = util.UniqueValues(append(p.lists, tag))
}

func (p *Payload) UpdateList(tag string, value []string) {
	keywords, ok := p.tags[tag].(string)
	if !ok {
		keywords = ""
	}

	values := strings.Split(keywords, separator)
	values = append(values, value...)
	p.AddList(tag, values)
}

func (p *Payload) AddTag(tag string, value string) {
	p.tags[tag] = filter(value)
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
