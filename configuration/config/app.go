package config

import (
	"github.com/vvval/go-metadata-scanner/configuration/vars"
	"gopkg.in/yaml.v2"
)

type App struct {
	toolPath   string
	extensions []string
	fields     []string
}

type AppSchema struct {
	ToolPath   string   `yaml:"exiftool"`
	Extensions []string `yaml:"extensions"`
	Fields     []string `yaml:"fields"`
}

func (c App) MergeDefaults() vars.Config {
	if len(c.toolPath) == 0 {
		c.toolPath = "exiftool"
	}

	return c
}

func (c App) Filename() string {
	return "./app.yaml"
}

func (c App) Schema() vars.Schema {
	return AppSchema{}
}

func (s AppSchema) Parse(data []byte) (vars.Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return App{}, err
	}

	return App{
		toolPath:   s.ToolPath,
		extensions: s.Extensions,
		fields:     s.Fields,
	}, nil
}

func (c App) ToolPath() string {
	return c.toolPath
}

func (c App) Extensions() []string {
	return c.extensions
}

func (c App) Fields() []string {
	return c.fields
}
