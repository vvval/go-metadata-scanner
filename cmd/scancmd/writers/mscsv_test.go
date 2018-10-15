package writers

import (
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestTitle(t *testing.T) {
	type check struct {
		f   vars.File
		exp string
	}

	set := []check{
		{
			vars.NewFile("file1.jpg", metadata.Tags{"XMP:Headline": "title test"}),
			"title test",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"title": "title test"}),
			"",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"XMP:Test": "title test"}),
			"",
		},
	}

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	for i, s := range set {
		v := title(&s.f, dict)
		if v != s.exp {
			t.Errorf("title mapping (line `%d`):\ngot `%s`\nexpected `%s`", i, v, s.exp)
		}
	}
}

func TestDescription(t *testing.T) {
	type check struct {
		f   vars.File
		exp string
	}

	set := []check{
		{
			vars.NewFile("file1.jpg", metadata.Tags{"EXIF:ImageDescription": "description test"}),
			"description test",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"description": "description test"}),
			"",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"XMP:ImageDescription": "title test"}),
			"",
		},
	}

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	for i, s := range set {
		v := description(&s.f, dict)
		if v != s.exp {
			t.Errorf("description mapping (line `%d`):\ngot `%s`\nexpected `%s`", i, v, s.exp)
		}
	}
}

func TestKeywords(t *testing.T) {
	type check struct {
		f   vars.File
		exp string
	}

	set := []check{
		{
			vars.NewFile("file1.jpg", metadata.Tags{"IPTC:Keywords": "keyword1"}),
			"keyword1",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"keywords": "keyword2"}),
			"",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"XMP:Subject": []string{"keyword1", "keyword2"}}),
			"keyword1, keyword2",
		},
		{
			vars.NewFile("file1.jpg", metadata.Tags{"XMP:Subject": 123}),
			"123",
		},
	}

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	for i, s := range set {
		v := keywords(&s.f, dict)
		if !reflect.DeepEqual(v, s.exp) {
			t.Errorf("keywords mapping (line `%d`):\ngot `%s`\nexpected `%s`", i, v, s.exp)
		}
	}
}
