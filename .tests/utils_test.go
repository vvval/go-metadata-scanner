package _tests

import (
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/rand"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"testing"
)

func TestTokenizer(t *testing.T) {
	set := map[string][]string{
		"a,b,,c, d": {"a", "b", "c", "d"},
		"a,b,a":     {"a", "b", "a"},
		`a,"b,c",d`: {"a", `b,c`, "d"},
		"a;b,c":     {"a", "b", "c"},
	}

	for str, exp := range set {
		tokens := util.SplitKeywords(str)
		if !reflect.DeepEqual(exp, tokens) {
			t.Errorf("tokens not equal:\ngot `%s`\nexpected `%s`", tokens, exp)
		}
	}
}

func TestAdjustSize(t *testing.T) {
	set := [][]int{
		{10, 3, 5, 2, 5},
		{10, 6, 5, 2, 5},
		{10, 3, 3, 3, 4},
		{10, 3, 2, 3, 4},
		{10, 3, 11, 1, 11},
		{10, 3, 9, 2, 9},
	}

	for _, v := range set {
		p, c := util.AdjustSizes(v[0], v[1], v[2])
		if p != v[3] || c != v[4] {
			t.Errorf("values are not equal:\ninput `%d`, `%d` and `%d`\ngot `%d` and `%d`\nexpected `%d` and `%d`", v[0], v[1], v[2], p, c, v[3], v[4])
		}
	}
}

func TestExtension(t *testing.T) {
	set := [][]string{
		{"filename.ext", "ext"},
		{".ext", "ext"},
		{"filename", ""},
		{"filename.", ""},
	}

	for _, str := range set {
		ext := util.Extension(str[0])
		if ext != str[1] {
			t.Errorf("extensions not equal:\ngot `%s`\nexpected `%s`", ext, str[1])
		}
	}
}

func TestRand(t *testing.T) {
	set := []int{10, 100}
	var reg = &regexp.Regexp{}
	reg = regexp.MustCompile("^[a-zA-Z0-9]*$")
	for _, n := range set {
		str := rand.Strings(n)
		if len(str) != n || !reg.MatchString(str) {
			t.Errorf("random string incorrect:\ngot `%s` of length `%d`\nexpected a-zA-A0-9 regex of `%d` length", str, len(str), n)
		}
	}
}

//func TestCandidates(t *testing.T) {
//	//t.Error("test file")
//}
//

func TestScanDir(t *testing.T) {
	type dir struct {
		dir      string
		ext      []string
		expected vars.Chunk
	}

	set := []dir{
		{"./assets", []string{"ext", "ext3"}, vars.Chunk{
			filepath.Join("assets", "subfolder1", "file1.ext"),
			filepath.Join("assets", "subfolder2", "file3.ext3"),
			filepath.Join("assets", "file2.ext"),
			filepath.Join("assets", "file4.ext3"),
		}},
		{"./assets/SubFolder2", []string{"ext2"}, vars.Chunk{}},
		{"./assets/SubFolder2", []string{"ext3"}, vars.Chunk{
			filepath.Join("assets", "subfolder2", "file3.ext3"),
		}},
		{"./assets/SubFolder2", []string{"ext3"}, vars.Chunk{
			filepath.Join("assets", "subfolder2", "file3.ext3"),
		}},
		{"./assets/SubFolder2", []string{}, vars.Chunk{
			filepath.Join("assets", "subfolder2", "file3.ext3"),
		}},
		{"./assets", []string{"ext2"}, vars.Chunk{
			filepath.Join("assets", "subfolder1", "subfolder3", "file5.ext2"),
		}},
	}

	for _, v := range set {
		res := scan.MustDir(v.dir, v.ext)
		sort.Strings(res)
		exp := v.expected
		sort.Strings(exp)
		if len(res) > 0 || len(exp) > 0 {
			if !reflect.DeepEqual(res, exp) {
				t.Errorf("scan dir incorrect for dir `%s` and ext %+v:\ngot `%+v`\nexpected `%+v`", v.dir, v.ext, res, exp)
			}
		}
	}
}
