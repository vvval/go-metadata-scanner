package configuration

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"io/ioutil"
)

func Load(conf Config, filename string) Config {
	var err error

	if util.FileExists(filename) {
		data, err := ioutil.ReadFile(filename)
		if err == nil {
			loadedConf, err := conf.Schema().Parse(data)
			if err == nil {
				return loadedConf.MergeDefault(conf)
			}
		}
	}

	if err != nil {
		log.Failure(fmt.Sprintf("Config \"%s\" read", filename), err.Error())
	}

	return conf
}
