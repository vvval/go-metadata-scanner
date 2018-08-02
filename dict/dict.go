package dict

import (
	"github.com/imdario/mergo"
	"github.com/vvval/go-metadata-scanner/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var conf config

const configFilename string = "./dict.yaml"

func Get() config {
	return conf
}

func init() {
	conf = load()
}

func load() config {
	defaultConfig := config{}

	if fileDetected() {
		fileConfig, err := loadFile()
		if err == nil {
			return mergeConfigs(fileConfig, defaultConfig)
		}

		log.Failure("Dict config read", err.Error())
	}

	return defaultConfig
}

func fileDetected() bool {
	_, err := os.Stat(configFilename)

	return err == nil
}

// Read file into config
func loadFile() (config, error) {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return config{}, err
	}

	//empty struct to fill data into
	conf := struct {
		Known    map[string][]string
		Booleans []string
		Lists    []string
	}{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return config{}, err
	}

	return config{conf.Known, conf.Booleans, conf.Lists}, nil
}

// Merge default config and file config.
// File config is primary
func mergeConfigs(fileConfig, defaultConfig config) config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
