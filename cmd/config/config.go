package config

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const configName string = "./config.yaml"
const exifToolPath string = "exiftool"

type config struct {
	ExifToolPath string `yaml:"exiftool"`
	Extensions   []string
	Fields       []string
	TagMap       map[string][]string
	BoolTags     []string
	ListTags     []string
}

var appConfig config

func AppConfig() config {
	return appConfig
}

func init() {
	appConfig = defineConfig()
}

func defineConfig() config {
	var defaultConfig = config{
		ExifToolPath: exifToolPath,
	}

	if configFileDetected() {
		fileConfig, err := loadConfig()
		if err != nil {
			log.Fatalln(err)
		}

		return mergeConfigs(fileConfig, defaultConfig)
	}

	return defaultConfig
}

func configFileDetected() bool {
	_, err := os.Stat(configName)

	return err == nil
}

// Read config file
func loadConfig() (config, error) {
	data, err := ioutil.ReadFile(configName)
	if err != nil {
		return config{}, err
	}

	fileConfig := config{}
	err = yaml.Unmarshal([]byte(data), &fileConfig)
	if err != nil {
		return config{}, err
	}

	return fileConfig, nil
}

// Merge default config and file config.
// File config is primary
func mergeConfigs(fileConfig, defaultConfig config) config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
