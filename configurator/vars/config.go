package vars

type Config interface {
	Filename() string
	Schema() Schema
}

type Schema interface {
	Parse(data []byte) (Config, error)
}

//var (
//	dict Config = struct {
//		known    map[string][]string
//		booleans []string
//		lists    []string
//	}{}
//	App Config = struct {
//	}{}
//	Mscsv Config = struct {
//	}{}
//)
