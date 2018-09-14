package config

import (
	"github.com/imdario/mergo"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const configFilename string = "./config.yaml"
const toolPath string = "exiftool"

var conf Config

func Get() Config {
	return conf
}

func init() {
	conf = load(Config{
		toolPath: toolPath,
	})
}

func load(defaultConfig Config) Config {
	if util.FileExists(configFilename) {
		fileConfig, err := loadFile()
		if err == nil {
			return mergeConfigs(fileConfig, defaultConfig)
		}

		log.Failure("App config read", err.Error())
	}

	return defaultConfig
}

// Read file into config
func loadFile() (Config, error) {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return Config{}, err
	}

	//empty struct to fill data into
	conf := struct {
		ToolPath   string   `yaml:"exiftool"`
		Extensions []string `yaml:"extensions"`
		Fields     []string `yaml:"fields"`
	}{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return Config{}, err
	}

	return Config{conf.ToolPath, conf.Extensions, conf.Fields}, nil
}

// Merge default config and file config.
// File config is primary
func mergeConfigs(fileConfig, defaultConfig Config) Config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
