package configuration

import (
	"github.com/vvval/go-metadata-scanner/util"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func TestAppLoad(t *testing.T) {
	type check struct {
		conf testConfig
		file string
		exp  testConfig
	}

	set := []check{
		{testConfig{"s", []string{"a1", "a2"}}, "", testConfig{"s", []string{"a1", "a2"}}},
		{testConfig{"s", []string{"a3", "a4"}}, "./fixtures/test4.yaml", testConfig{"s", []string{"a3", "a4"}}},
		{testConfig{"s", []string{}}, "./fixtures/test1.yaml", testConfig{"str", []string{"arr1", "arr2"}}},
		{testConfig{"s", []string{"a5", "a6"}}, "./fixtures/test1.yaml", testConfig{"str", []string{"arr1", "arr2", "a5", "a6"}}},
		{testConfig{"s", []string{"a7", "a8"}}, "./fixtures/test2.yaml", testConfig{"str", []string{"a7", "a8"}}},
		{testConfig{"s", []string{"a9", "a10"}}, "./fixtures/test3.yaml", testConfig{"s", []string{"arr1", "arr2", "a9", "a10"}}},
	}

	for i, s := range set {
		l := Load(s.conf, s.file)
		if !reflect.DeepEqual(l, s.exp) {
			t.Errorf("load failed (line `%d`):\ngot `%s`\nexpected `%s`", i, l, s.exp)
		}
	}
}

type testSchema struct {
	StringValue string   `yaml:"string"`
	ArrayValue  []string `yaml:"array"`
}

func (c testConfig) Schema() Schema {
	return testSchema{}
}

func (s testSchema) Parse(data []byte) (Config, error) {
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return testConfig{}, err
	}

	return testConfig{s.StringValue, s.ArrayValue}, nil
}

func (c testConfig) MergeDefault(conf Config) Config {
	if len(c.stringValue) == 0 {
		c.stringValue = conf.(testConfig).stringValue
	}

	c.arrayValue = util.UniqueValues(append(c.arrayValue, conf.(testConfig).arrayValue...))

	return c
}

type testConfig struct {
	stringValue string
	arrayValue  []string
}
