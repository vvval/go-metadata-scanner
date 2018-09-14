package configurator

import "github.com/vvval/go-metadata-scanner/configurator/config"

var (
	Dict config.Dict
	//app  config.Dict
)

func LoadDict() {
	Load(Dict)
}
