package config

type config struct {
	toolPath   string
	extensions []string
	fields     []string
}

func (c config) ToolPath() string {
	return c.toolPath
}

func (c config) Extensions() []string {
	return c.extensions
}

func (c config) Fields() []string {
	return c.fields
}
