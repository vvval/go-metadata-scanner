package scan

import (
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"path/filepath"
	"testing"
)

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
		c, _ := Candidates(v.file, v.files, v.ext)
		if c != v.exp {
			t.Errorf("candidates incorrect for file `%s` and ext %+v:\ngot `%+v`\nexpected `%+v`", v.file, v.ext, c, v.exp)
		}
	}
}

func TestScanDir(t *testing.T) {
	type check struct {
		dir string
		ext []string
		exp vars.Chunk
	}

	set := []check{
		{"./fixtures", []string{"ext", "ext3"}, vars.Chunk{
			filepath.Join("fixtures", "subfolder1", "file1.ext"),
			filepath.Join("fixtures", "subfolder2", "file3.ext3"),
			filepath.Join("fixtures", "file2.ext"),
			filepath.Join("fixtures", "file4.ext3"),
		}},
		{"./fixtures/SubFolder2", []string{"ext2"}, vars.Chunk{}},
		{"./fixtures/SubFolder2", []string{"ext3"}, vars.Chunk{
			filepath.Join("fixtures", "subfolder2", "file3.ext3"),
		}},
		{"./fixtures/SubFolder2", []string{"ext3"}, vars.Chunk{
			filepath.Join("fixtures", "subfolder2", "file3.ext3"),
		}},
		{"./fixtures/SubFolder2", []string{}, vars.Chunk{
			filepath.Join("fixtures", "subfolder2", "file3.ext3"),
		}},
		{"./fixtures", []string{"ext2"}, vars.Chunk{
			filepath.Join("fixtures", "subfolder1", "subfolder3", "file5.ext2"),
		}},
	}

	for _, v := range set {
		res := MustDir(v.dir, v.ext)
		exp := v.exp
		if !util.Equal(res, exp) && (len(res) > 0 || len(exp) > 0) {
			t.Errorf("scan dir incorrect for dir `%s` and ext %+v:\ngot `%+v`\nexpected `%+v`", v.dir, v.ext, res, exp)
		}
	}
}
