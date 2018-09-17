package writers

import "fmt"

func Get(ext string) (Writer, error) {
	switch ext {
	case "csv":
		return &CSVWriter{}, nil

	case "mscsv":
		return &MSCSVWriter{}, nil

	case "json":
		return &JSONWriter{}, nil
	}

	return nil, fmt.Errorf("unsupported writer type `%s`", ext)
}
