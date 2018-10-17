package configuration

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"io/ioutil"
	"reflect"
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
	} else {
		log.Failure("Config load failed, file not exists", getType(conf))
	}

	if err != nil {
		log.Failure(fmt.Sprintf("Config \"%s\" read", filename), err.Error())
	}

	return conf
}

func getType(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
