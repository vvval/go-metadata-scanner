package cmd

import (
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const configName string = "./config.yaml"
const exifToolPath string = "./exiftool"

type config struct {
	ExifToolPath string `yaml:"exiftool"`
}

func defineConfig() config {
	var defaultConfig = config{exifToolPath}

	if configFileDetected() {
		fileConfig := loadConfig()

		return mergeConfigs(fileConfig, defaultConfig)
	}

	return defaultConfig
}

func configFileDetected() bool {
	_, err := os.Stat(configName)

	return err == nil
}

// Read config file, panic when reading or unmarshal fails
func loadConfig() config {
	data, err := ioutil.ReadFile(configName)
	if err != nil {
		panic(fmt.Sprintf("Error occurred while reading config %s: %s", configName, err))
	}

	fileConfig := config{}
	err = yaml.Unmarshal([]byte(data), &fileConfig)
	if err != nil {
		panic(fmt.Sprintf("Error occurred while unmarshalling config %s: %s", configName, err))
	}

	return fileConfig
}

func mergeConfigs(fileConfig, defaultConfig config) config {
	mergo.Merge(&fileConfig, defaultConfig)

	return fileConfig
}
