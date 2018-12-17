package config

import (
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"testing"
)

func TestAppConfig(t *testing.T) {
	type check struct {
		a, b, exp AppConfig
	}

	set := []check{
		{
			AppConfig{"path", []string{"ext1", "ext2"}, []string{"f1", "f2"}},
			AppConfig{"", []string{}, []string{}},
			AppConfig{"path", []string{"ext1", "ext2"}, []string{"f1", "f2"}},
		},
		{
			AppConfig{"", []string{}, []string{}},
			AppConfig{"path", []string{"ext1", "ext2"}, []string{"f1", "f2"}},
			AppConfig{"path", []string{"ext1", "ext2"}, []string{"f1", "f2"}},
		},
		{
			AppConfig{"path", []string{"ext1", "ext2"}, []string{"f1", "f2"}},
			AppConfig{"path2", []string{"ext3", "ext2"}, []string{"f3", "f2"}},
			AppConfig{"path", []string{"ext1", "ext2", "ext3"}, []string{"f1", "f2", "f3"}},
		},
	}

	for i, s := range set {
		m := s.a.MergeDefault(s.b).(AppConfig)
		if m.toolPath != s.exp.toolPath || !util.Equal(m.fields, s.exp.fields) || !util.Equal(m.extensions, s.exp.extensions) {
			t.Errorf("merge failed (line `%d`):\ngot `%v`\nexpected `%v`", i, m, s.exp)
		}
	}
}

func TestDictConfig(t *testing.T) {
	type check struct {
		a, b, exp DictConfig
	}

	set := []check{
		{
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, map[string]string{}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{}, map[string]string{}, []string{}, []string{}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, map[string]string{}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
		},
		{
			DictConfig{map[string][]string{}, map[string]string{}, []string{}, []string{}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, map[string]string{}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, map[string]string{}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
		},
		{
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, map[string]string{}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{"b": {"b1", "b2"}}, map[string]string{}, []string{"bool3", "bool2"}, []string{"list3", "list2"}},
			DictConfig{map[string][]string{"a": {"a2", "a1"}, "b": {"b1", "b2"}}, map[string]string{}, []string{"bool1", "bool2", "bool3"}, []string{"list1", "list2", "list3"}},
		},
	}

	for i, s := range set {
		m := s.a.MergeDefault(s.b).(DictConfig)
		if !mapEqual(m.known, s.exp.known) || !util.Equal(m.booleans, s.exp.booleans) || !util.Equal(m.lists, s.exp.lists) {
			t.Errorf("merge failed (line `%d`):\ngot `%v`\nexpected `%v`", i, m, s.exp)
		}
	}
}

func TestDictFind(t *testing.T) {
	type check struct {
		name  string
		found bool
	}

	set := []check{
		{"", false},
		{"test", false},
		{"title", true},
		{"IPTC:Headline", true},
		{"iptc:Headline", true},
		{"keywords", true},
		{"IPTC:keywords", true},
		{"iptc:KEYWORDS", true},
	}

	dict := configuration.Load(Dict, "./../dict.yaml").(DictConfig)

	for i, s := range set {
		_, f := dict.Find(s.name)
		if f != s.found {
			t.Errorf("find failed (line `%d`):\ngot `%t`\nexpected `%t`", i, f, s.found)
		}
	}
}

func TestDictFindGroup(t *testing.T) {
	type check struct {
		name  string
		found bool
	}

	set := []check{
		{"keywords", true},
		{"kkeywords", false},
		{"keywordss", true},
	}

	dict := configuration.Load(Dict, "./../dict.yaml").(DictConfig)

	for i, s := range set {
		_, f := dict.Find(s.name)
		if f != s.found {
			t.Errorf("find groups failed (line `%d`):\ngot `%t`\nexpected `%t`", i, f, s.found)
		}
	}
}

func TestDictIsBoolean(t *testing.T) {
	type check struct {
		name string
		is   bool
	}

	set := []check{
		{"", false},
		{"test", false},
		{"copyrighted", true},
		{"XMP:marked", true},
	}

	dict := configuration.Load(Dict, "./../dict.yaml").(DictConfig)

	for i, s := range set {
		tag, ok := dict.Find(s.name)
		if !ok && s.is {
			t.Errorf("find in booleans failed (line `%d`) for `%+v` (unknown tag)", i, s)
		}

		found := isFound(tag, dict)

		if found != s.is {
			t.Errorf("find in booleans failed (line `%d`) for `%+v`:\ngot `%t`\nexpected `%t`", i, s, found, s.is)
		}
	}
}

func isFound(tag vars.Tag, dict DictConfig) bool {
	for _, mapTag := range tag.Map() {
		if dict.IsBoolean(tag.Key(), mapTag) {
			return true
		}
	}

	return false
}

func TestDictIsList(t *testing.T) {
	type check struct {
		name, tag string
		is        bool
	}

	set := []check{
		{"", "", false},
		{"test", "", false},
		{"keywords", "test", true},
		{"", "IPTC:Writer-Editor", true},
		{"", "captionWriter", false},
	}

	dict := configuration.Load(Dict, "./../dict.yaml").(DictConfig)

	for i, s := range set {
		f := dict.IsList(s.name, s.tag)
		if f != s.is {
			t.Errorf("find in lists failed (line `%d`) for `%+v`:\ngot `%t`\nexpected `%t`", i, s, f, s.is)
		}
	}
}

func TestDictTagEquals(t *testing.T) {
	type check struct {
		t1, t2 string
		exp    bool
	}

	set := []check{
		{"", "", true},
		{"test", "", false},
		{"", "test", false},
		{"test", "test", true},
		{"IPTC:Contact", "IPTC:Contact", true},
		{"contact", "IPTC:Contact", true},
		{"IPTC:Contact", "contact", false},
	}

	for i, s := range set {
		f := tagEquals(s.t1, s.t2)
		if f != s.exp {
			t.Errorf("tags equality failed (line `%d`):\ngot `%t`\nexpected `%t`", i, f, s.exp)
		}
	}
}

func TestMSCSV(t *testing.T) {
	type check struct {
		a, b, exp MSCSVConfig
	}

	set := []check{
		{MSCSVConfig{"provider"}, MSCSVConfig{""}, MSCSVConfig{"provider"}},
		{MSCSVConfig{""}, MSCSVConfig{"provider"}, MSCSVConfig{"provider"}},
		{MSCSVConfig{"provider"}, MSCSVConfig{"provider2"}, MSCSVConfig{"provider"}},
	}

	for i, s := range set {
		m := s.a.MergeDefault(s.b).(MSCSVConfig)
		if m.provider != s.exp.provider {
			t.Errorf("merge failed (line `%d`):\ngot `%v`\nexpected `%v`", i, m, s.exp)
		}
	}
}

func mapEqual(a, b map[string][]string) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		vv, ok := b[k]
		if !ok || !util.Equal(v, vv) {
			return false
		}
	}

	return true
}
