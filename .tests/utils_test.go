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
	type check struct {
		n, d, min, an, ad int
	}
	set := []check{
		{10, 3, 5, 2, 5},
		{10, 6, 5, 2, 5},
		{10, 3, 3, 3, 4},
		{10, 3, 2, 3, 4},
		{10, 3, 11, 1, 11},
		{10, 3, 9, 2, 9},
	}

	for _, v := range set {
		p, c := util.AdjustSizes(v.n, v.d, v.min)
		if p != v.an || c != v.ad {
			t.Errorf("values are not equal:\ninput `%d`, `%d` and `%d`\ngot `%d` and `%d`\nexpected `%d` and `%d`", v.n, v.d, v.min, p, c, v.an, v.ad)
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

func TestCandidates(t *testing.T) {
	type check struct {
		file  string
		files vars.Chunk
		ext   []string
		exp   string
	}

	set := []check{
		// Adding extension
		{"1", vars.Chunk{"folder/1.png", "folder/2.png"}, []string{"png"}, "folder/1.png"},
		// Adding extension even if file has it
		{"1.png", vars.Chunk{"folder/1.png.ext", "folder/2.png"}, []string{"png", "ext"}, "folder/1.png.ext"},
		// Adding prefixes ([a-z]?_0?)
		{"1.png", vars.Chunk{"folder/01.png"}, []string{"png"}, "folder/01.png"},
		{"1.png", vars.Chunk{"folder/img_01.png"}, []string{"png"}, "folder/img_01.png"},
		// If prefix has letters, it should end with "underscore" char
		{"1.png", vars.Chunk{"folder/img01.png"}, []string{"png"}, ""},
		// Prefix can't have prefix made of non-zero digits
		{"1.png", vars.Chunk{"folder/21.png"}, []string{"png"}, ""},
		// If more than 1 candidate - ignore file
		{"1.png", vars.Chunk{"folder/01.png", "folder/001.png"}, []string{"png", "ext"}, ""},
		// But if has full match - it's fine
		{"folder2/1", vars.Chunk{"folder2/1.png", "folder/01.png", "folder/001.png"}, []string{"png", "ext"}, "folder2/1.png"},
		// This is not a full match, folder prefix is missing
		{"1", vars.Chunk{"folder2/1.png", "folder/01.png", "folder/001.png"}, []string{"png", "ext"}, ""},
	}

	for _, v := range set {
		c, _ := scan.Candidates(v.file, v.files, v.ext)
		if c != v.exp {
			t.Errorf("candidates incorrect for file `%s` and ext %+v:\ngot `%+v`\nexpected `%+v`", v.file, v.ext, c, v.exp)
		}
	}
}

func TestScanDir(t *testing.T) {
	type check struct {
		dir      string
		ext      []string
		expected vars.Chunk
	}

	set := []check{
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
		if !reflect.DeepEqual(res, exp) && (len(res) > 0 || len(exp) > 0) {
			t.Errorf("scan dir incorrect for dir `%s` and ext %+v:\ngot `%+v`\nexpected `%+v`", v.dir, v.ext, res, exp)
		}
	}
}
