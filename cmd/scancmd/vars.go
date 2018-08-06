package scancmd

var (
	Input    Flags
	Files    chan []string
	PoolSize = 10
)

type Flags struct {
	directory string
	format    string
	filename  string
}

func (input Flags) Filename() string {
	return input.filename
}

func (input Flags) Directory() string {
	return input.directory
}

func (input Flags) Format() string {
	return input.format
}
