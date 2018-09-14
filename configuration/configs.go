package configuration

import "github.com/vvval/go-metadata-scanner/configuration/config"

var (
	Dict config.Dict
	App  config.App
)

func LoadAll() {
	Dict = Load(Dict).(config.Dict)
	App = Load(App).(config.App)
}
