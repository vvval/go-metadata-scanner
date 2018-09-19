package config

import (
	"github.com/vvval/go-metadata-scanner/util"
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
			t.Errorf("merge failed for %d:\ngot %v\nexpected %v", i, m, s.exp)
		}
	}
}

func TestDictConfig(t *testing.T) {
	type check struct {
		a, b, exp DictConfig
	}

	set := []check{
		{
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{}, []string{}, []string{}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
		},
		{
			DictConfig{map[string][]string{}, []string{}, []string{}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
		},
		{
			DictConfig{map[string][]string{"a": {"a1", "a2"}}, []string{"bool1", "bool2"}, []string{"list1", "list2"}},
			DictConfig{map[string][]string{"b": {"b1", "b2"}}, []string{"bool3", "bool2"}, []string{"list3", "list2"}},
			DictConfig{map[string][]string{"a": {"a2", "a1"}, "b": {"b1", "b2"}}, []string{"bool1", "bool2", "bool3"}, []string{"list1", "list2", "list3"}},
		},
	}

	for i, s := range set {
		m := s.a.MergeDefault(s.b).(DictConfig)
		if !mapEqual(m.known, s.exp.known) || !util.Equal(m.booleans, s.exp.booleans) || !util.Equal(m.lists, s.exp.lists) {
			t.Errorf("merge failed for %d:\ngot %v\nexpected %v", i, m, s.exp)
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
			t.Errorf("merge failed for %d:\ngot %v\nexpected %v", i, m, s.exp)
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
