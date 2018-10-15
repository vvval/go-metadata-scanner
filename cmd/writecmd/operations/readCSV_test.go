package operations

import (
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	filename := "./fixtures/file1.csv"
	file := util.MustOpenReadonlyFile(filename)
	defer file.Close()

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	type check struct {
		filename string
		lists    []string
		tags     []string
		miss     []string
	}

	expected := map[string]check{
		"file.jpg":  {"file.jpg", []string{}, []string{"IPTC:Headline", "IPTC:Keywords", "XMP:Marked"}, []string{"EXIF:ImageDescription", "IPTC:Caption-Abstract"}},
		"file2.jpg": {"file2.jpg", []string{"IPTC:Keywords"}, []string{"XMP:Headline", "XMP:Subject"}, []string{"XMP:Description", "some other description"}},
	}
	var output []check

	ReadCSV(util.GetCSVReader(file, ','), dict, func(filename string, payload metadata.Payload) {
		var tags, miss []string
		for s, v := range payload.Tags() {
			tags = append(tags, s)
			if vs, ok := v.(string); ok {
				if len(vs) == 0 {
					miss = append(miss, s)
				}
			} else if v != nil {
				miss = append(miss, s)
			}
		}

		output = append(output, check{filename, payload.ListTags(), tags, miss})
	})

	for i, o := range output {
		e, ok := expected[o.filename]
		if !ok {
			t.Errorf("read unexpected result (line `%d`) for filename `%s`", i, o.filename)
		}

		for _, l := range e.lists {
			if !inArray(l, o.lists) {
				t.Errorf("read expected lists not found (line `%d`) for filename `%s`:\ntag %s\ngot: %+v", i, o.filename, l, o.lists)
			}
		}
	}
}

func inArray(el string, arr []string) bool {
	for _, val := range arr {
		if strings.EqualFold(el, val) {
			return true
		}
	}

	return false
}
