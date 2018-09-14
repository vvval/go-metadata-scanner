package dict

import (
	"github.com/imdario/mergo"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const configFilename string = "./dict.yaml"

var conf Config

func Get() Config {
	return conf
}

func init() {
	conf = load(Config{})
}

func load(defaultConfig Config) Config {
	if util.FileExists(configFilename) {
		fileConfig, err := loadFile()
		if err == nil {
			return mergeConfigs(fileConfig, defaultConfig)
		}

		log.Failure("Dict config read", err.Error())
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
		Known    map[string][]string `yaml:"known"`
		Booleans []string            `yaml:"booleans"`
		Lists    []string            `yaml:"lists"`
	}{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return Config{}, err
	}

	return Config{conf.Known, conf.Booleans, conf.Lists}, nil
}

// Merge default config and file config.
// File config is primary
func mergeConfigs(fileConfig, defaultConfig Config) Config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
