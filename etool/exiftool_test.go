package etool

import (
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"reflect"
	"testing"
)

func TestReadArgs(t *testing.T) {
	type check struct {
		names  []string
		fields []string
		exp    []string
	}

	set := []check{
		{[]string{"a", "b"}, []string{"f1", "f2"}, []string{"-j", "-G", "-b", "-f1:all", "-f2:all", "a", "b"}},
		{[]string{}, []string{"f1", "f2"}, []string{"-j", "-G", "-b", "-f1:all", "-f2:all"}},
		{[]string{"a", "b"}, []string{}, []string{"-j", "-G", "-b", "a", "b"}},
	}

	for i, v := range set {
		p := packReadArgs(v.names, v.fields)
		if !reflect.DeepEqual(p, v.exp) {
			t.Errorf("args mismatch (line `%d`):\ninput `%s` and `%s`:\ngot `%s`\nexpected `%s`", i, v.names, v.fields, p, v.exp)
		}
	}
}

func TestWriteArgs(t *testing.T) {
	type check struct {
		name      string
		tags      metadata.Tags
		useSep    bool
		originals bool
		exp       []string
	}

	set := []check{
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, true, true,
			[]string{"-n1=v1", "-n2=v2,v3", "-sep", metadata.Separator(), "name"}},
		{"", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, true, true,
			[]string{"-n1=v1", "-n2=v2,v3", "-sep", metadata.Separator(), ""}},
		{"", metadata.Tags{}, true, true,
			[]string{"-sep", metadata.Separator(), ""}},
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, false, true,
			[]string{"-n1=v1", "-n2=v2,v3", "name"}},
		{"name", metadata.Tags{"n1": "v1", "n2": "v2,v3"}, false, false,
			[]string{"-n1=v1", "-n2=v2,v3", overwriteFlag, "name"}},
	}
	for i, v := range set {
		p := packWriteArgs(v.name, v.tags, v.useSep, v.originals)
		if !util.Equal(p, v.exp) {
			t.Errorf("args mismatch (line `%d`) check\ngot `%+v`\nexp `%+v`", i, p, v.exp)
		}
	}
}

func TestParse(t *testing.T) {
	type check struct {
		data string
		exp  []vars.File
	}

	set := []check{
		{`[{"SourceFile": "folder/test1.jpg","File:FileName": "test1.jpg"},{"SourceFile": "folder/test2.jpg","File:FileName": "test2.jpg"}]`, []vars.File{
			vars.NewFile("folder/test1.jpg", metadata.Tags{"File:FileName": "test1.jpg"}),
			vars.NewFile("folder/test2.jpg", metadata.Tags{"File:FileName": "test2.jpg"}),
		}},
		{`[]`, []vars.File{}},
		{`[{"SourceFile": "folder/test1.jpg","File:FileName": "test1.jpg"},{"SourceFile2": "folder/test2.jpg","File:FileName": "test2.jpg"}]`, []vars.File{
			vars.NewFile("folder/test1.jpg", metadata.Tags{"File:FileName": "test1.jpg"}),
		}},
		{`[{"SourceFile": "folder/test1.jpg","File:FileName": "test1.jpg"},{"File:FileName": "test2.jpg"}]`, []vars.File{
			vars.NewFile("folder/test1.jpg", metadata.Tags{"File:FileName": "test1.jpg"}),
		}},
	}

	for i, s := range set {
		p := Parse([]byte(s.data))
		if !reflect.DeepEqual(p, s.exp) && len(p) > 0 && len(s.exp) > 0 {
			t.Errorf("parse failed (line `%d`) check\ngot `%+v`\nexp `%+v`", i, p, s.exp)
		}
	}
}
