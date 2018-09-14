package writers

import "fmt"

func Get(ext string, filename string, headers []string) (Writer, error) {
	switch ext {
	case "csv":
		return NewCSVWriter(filename, headers), nil

	case "json":
		return NewJSONWriter(filename, headers), nil
	}

	return nil, fmt.Errorf("unsupported writer type `%s`", ext)
}
