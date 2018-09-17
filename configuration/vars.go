package configuration

type Config interface {
	Schema() Schema
	MergeDefault(conf Config) Config
}

type Schema interface {
	Parse(data []byte) (Config, error)
}
