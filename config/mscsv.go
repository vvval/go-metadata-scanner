package config

import (
	"github.com/vvval/go-metadata-scanner/configuration"
	"gopkg.in/yaml.v2"
)

type MSCSVSchema struct {
	Provider string `yaml:"provider"`
}

func (c MSCSVConfig) Schema() configuration.Schema {
	return MSCSVSchema{}
}

func (s MSCSVSchema) Parse(data []byte) (configuration.Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return MSCSVConfig{}, err
	}

	return MSCSVConfig{provider: s.Provider}, nil
}

func (c MSCSVConfig) MergeDefault(conf configuration.Config) configuration.Config {
	if len(c.provider) == 0 {
		c.provider = conf.(MSCSVConfig).provider
	}

	return c
}

type MSCSVConfig struct {
	provider string
}

func (c MSCSVConfig) Provider() string {
	return c.provider
}
