package config

import (
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"gopkg.in/yaml.v2"
	"strings"
)

type DictSchema struct {
	Known    map[string][]string `yaml:"known"`
	Booleans []string            `yaml:"booleans"`
	Lists    []string            `yaml:"lists"`
}

func (c DictConfig) Schema() configuration.Schema {
	return DictSchema{}
}

func (s DictSchema) Parse(data []byte) (configuration.Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return DictConfig{}, err
	}

	return DictConfig{
		known:    s.Known,
		booleans: s.Booleans,
		lists:    s.Lists,
	}, nil
}

func (c DictConfig) MergeDefault(conf configuration.Config) configuration.Config {
	for k, v := range conf.(DictConfig).known {
		if _, ok := c.known[k]; !ok {
			c.known[k] = v
		}
	}

	c.booleans = util.UniqueValues(append(c.booleans, conf.(DictConfig).booleans...))
	c.lists = util.UniqueValues(append(c.lists, conf.(DictConfig).lists...))

	return c
}

type DictConfig struct {
	known    map[string][]string
	booleans []string
	lists    []string
}

func (c DictConfig) Find(name string) (vars.Tag, bool) {
	if tag, found := known(name, c.known); found {
		return tag, found
	}

	for _, b := range c.booleans {
		if strings.EqualFold(b, name) {
			return found("", name, []string{name})
		}
	}

	for _, l := range c.lists {
		if strings.EqualFold(l, name) {
			return found("", name, []string{name})
		}
	}

	return notFound(name)
}

func (c DictConfig) IsBoolean(key, tag string) bool {
	return oneOf(key, tag, c.booleans)
}

func (c DictConfig) IsList(key, tag string) bool {
	return oneOf(key, tag, c.lists)
}

func known(name string, lists map[string][]string) (vars.Tag, bool) {
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

func found(key, name string, list []string) (vars.Tag, bool) {
	return vars.NewFoundTag(key, name, list), true
}

func notFound(name string) (vars.Tag, bool) {
	return vars.NewNotFoundTag(name), false
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
