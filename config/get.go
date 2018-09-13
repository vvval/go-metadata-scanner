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

var conf config

func Get() config {
	return conf
}

func init() {
	conf = load(config{
		toolPath: toolPath,
	})
}

func load(defaultConfig config) config {
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
func loadFile() (config, error) {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return config{}, err
	}

	//empty struct to fill data into
	conf := struct {
		ToolPath   string   `yaml:"exiftool"`
		Extensions []string `yaml:"extensions"`
		Fields     []string `yaml:"fields"`
	}{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return config{}, err
	}

	return config{conf.ToolPath, conf.Extensions, conf.Fields}, nil
}

// Merge default config and file config.
// File config is primary
func mergeConfigs(fileConfig, defaultConfig config) config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
