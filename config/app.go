package config

import (
	"github.com/vvval/go-metadata-scanner/configuration"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	toolPath   string
	extensions []string
	fields     []string
}

type AppSchema struct {
	ToolPath   string   `yaml:"exiftool"`
	Extensions []string `yaml:"extensions"`
	Fields     []string `yaml:"fields"`
}

func (c AppConfig) MergeDefault(conf configuration.Config) configuration.Config {
	if len(c.toolPath) == 0 {
		if conf, ok := conf.(AppConfig); ok {
			c.toolPath = conf.toolPath
		}
	}

	return c
}

func (c AppConfig) Schema() configuration.Schema {
	return AppSchema{}
}

func (s AppSchema) Parse(data []byte) (configuration.Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		toolPath:   s.ToolPath,
		extensions: s.Extensions,
		fields:     s.Fields,
	}, nil
}

func (c AppConfig) ToolPath() string {
	return c.toolPath
}

func (c AppConfig) Extensions() []string {
	return c.extensions
}

func (c AppConfig) Fields() []string {
	return c.fields
}
