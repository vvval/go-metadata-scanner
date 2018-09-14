package configuration

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/configuration/vars"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"io/ioutil"
)

func Load(conf vars.Config) vars.Config {
	var err error

	if util.FileExists(conf.Filename()) {
		data, err := ioutil.ReadFile(conf.Filename())
		if err == nil {
			loadedConf, err := conf.Schema().Parse(data)
			if err == nil {
				return loadedConf.MergeDefaults()
			}
		}
	}

	if err != nil {
		log.Failure(fmt.Sprintf("Config \"%s\" read", conf.Filename()), err.Error())
	}

	return conf.MergeDefaults()
}
