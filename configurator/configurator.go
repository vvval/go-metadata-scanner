package configurator

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/vvval/go-metadata-scanner/configurator/vars"
	"github.com/vvval/go-metadata-scanner/util/log"
	"io/ioutil"
	"os"
)

func Load(conf vars.Config) vars.Config {
	if fileExists(conf.Filename()) {
		data, err := ioutil.ReadFile(conf.Filename())
		if err != nil {
			log.Failure(fmt.Sprintf("Config \"%s\" read", conf.Filename()), err.Error())

			return conf
		}

		loadedConfig, err := conf.Schema().Parse(data)
		if err != nil {
			log.Failure(fmt.Sprintf("Config \"%s\" read", conf.Filename()), err.Error())

			return conf
		}

		return mergeConfigs(loadedConfig, conf)
	}

	return conf
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}

// Merge defaulted config and loaded from file config.
// File config is primary
func mergeConfigs(loaded, defaulted vars.Config) vars.Config {
	mergo.Merge(&loaded, defaulted)

	return loaded
}
