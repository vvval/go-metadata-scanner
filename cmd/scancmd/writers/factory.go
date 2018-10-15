package writers

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
)

func Get(ext string) (Writer, error) {
	switch ext {
	case "csv":
		return &CSVWriter{}, nil

	case "mscsv":
		return &MSCSVWriter{dict: config.Dict, config: config.MSCSV}, nil

	case "json":
		return &JSONWriter{}, nil
	}

	return nil, fmt.Errorf("unsupported writer type `%s`", ext)
}
