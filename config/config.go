package config

type Config struct {
	toolPath   string
	extensions []string
	fields     []string
}

func (c Config) ToolPath() string {
	return c.toolPath
}

func (c Config) Extensions() []string {
	return c.extensions
}

func (c Config) Fields() []string {
	return c.fields
}
