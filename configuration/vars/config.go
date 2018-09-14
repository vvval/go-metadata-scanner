package vars

type Config interface {
	Filename() string
	Schema() Schema
	MergeDefaults() Config
}

type Schema interface {
	Parse(data []byte) (Config, error)
}
