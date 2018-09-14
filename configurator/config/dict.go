package config

import (
	"github.com/vvval/go-metadata-scanner/configurator/vars"
	"github.com/vvval/go-metadata-scanner/dict"
	"gopkg.in/yaml.v2"
	"strings"
)

const configFilename string = "dict.yaml"

type Dict struct {
	known    map[string][]string
	booleans []string
	lists    []string
}

type DictSchema struct {
	Known    map[string][]string `yaml:"known"`
	Booleans []string            `yaml:"booleans"`
	Lists    []string            `yaml:"lists"`
}

func (conf Dict) Filename() string {
	return configFilename
}

func (conf Dict) Schema() vars.Schema {
	return DictSchema{}
}

func (s DictSchema) Parse(data []byte) (vars.Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return Dict{}, err
	}

	return Dict{
		known:    s.Known,
		booleans: s.Booleans,
		lists:    s.Lists,
	}, nil
}

func (conf Dict) Find(name string) (dict.Tag, bool) {
	if tag, found := known(name, conf.known); found {
		return tag, found
	}

	for _, b := range conf.booleans {
		if strings.EqualFold(b, name) {
			return found("", name, []string{name})
		}
	}

	for _, l := range conf.lists {
		if strings.EqualFold(l, name) {
			return found("", name, []string{name})
		}
	}

	return notFound(name)
}

func (conf Dict) IsBoolean(key, tag string) bool {
	return oneOf(key, tag, conf.booleans)
}

func (conf Dict) IsList(key, tag string) bool {
	return oneOf(key, tag, conf.lists)
}

func known(name string, lists map[string][]string) (dict.Tag, bool) {
	for key, list := range lists {
		if strings.EqualFold(key, name) {
			return found(key, name, list)
		}

		for _, val := range list {
			if tagEquals(val, name) {
				return found(key, name, list)
			}
		}
	}

	return notFound(name)
}

func found(key, name string, list []string) (dict.Tag, bool) {
	return dict.NewFoundTag(key, name, list), true
}

func notFound(name string) (dict.Tag, bool) {
	return dict.NewNotFoundTag(name), false
}

func oneOf(key, tag string, set []string) bool {
	for _, val := range set {
		if strings.EqualFold(tag, val) || strings.EqualFold(key, val) {
			return true
		}
	}

	return false
}

// Name or truncated name (without "<GROUP>:" prefix) equals
// Tag argument is a full tag name
// Name argument is a searchable input+
func tagEquals(tag, name string) bool {
	return strings.EqualFold(tag, name) || strings.EqualFold(tag, truncatePrefix(name))
}

// Cut "<group>:" prefix if found
func truncatePrefix(tag string) string {
	prefixEnding := strings.Index(tag, ":")
	if prefixEnding == -1 {
		return tag
	}

	runes := []rune(tag)

	return string(runes[prefixEnding+1:])
}
