package writers

import (
	"encoding/csv"
	"fmt"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type MSCSVWriter struct {
	BaseWriter
	csv *csv.Writer
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *MSCSVWriter) Write(file *vars.File) error {
	record, err := packImage(file)
	if err != nil {
		return err
	}

	return w.csv.Write(record)
}

func packImage(f *vars.File) ([]string, error) {
	file, err := os.Open(f.Filename())
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	record := []string{
		config.MSCSV.Provider(),
		f.RelPath(),
		fmt.Sprintf("%d", stat.Size()),
		mime.TypeByExtension(filepath.Ext(file.Name())),
		description(f),
		util.Extension(f.Filename()),
		fmt.Sprintf("%d", img.Width),
		fmt.Sprintf("%d", img.Height),
		keywords(f),
		title(f),
	}

	return record, nil
}

func (w *MSCSVWriter) Open(filename string, h []string) error {
	w.BaseWriter = NewWriter(filename, headers())

	file, err := openFile(w.filename)
	if err != nil {
		return err
	}

	w.file = file
	w.csv = csv.NewWriter(file)
	w.csv.Write(headers())

	return nil
}

func (w *MSCSVWriter) Close() error {
	if w.csv != nil {
		w.csv.Flush()
	}

	closeFile(w.file)

	return nil
}

func description(file *vars.File) string {
	return findTagValue(file, "description").(string)
}

func keywords(file *vars.File) string {
	keywords := findTagValue(file, "keywords")

	if str, ok := keywords.(string); ok {
		return str
	}

	value := reflect.ValueOf(keywords)
	output := make([]string, value.Len())

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			output[i] = fmt.Sprintf("%s", value.Index(i))
		}
		return strings.Join(output, ", ")
	}

	return fmt.Sprintf("%s", keywords)
}

func title(file *vars.File) string {
	return findTagValue(file, "title").(string)
}

func findTagValue(file *vars.File, field string) interface{} {
	if tag, found := config.Dict.Find(field); found {
		for _, name := range tag.Map() {
			if value, ok := file.Tags()[name]; ok && value != nil {
				return value
			}
		}
	}

	return ""
}

func headers() []string {
	return []string{
		"Provider",       //ELearningBrothers
		"MediaId",        //File name, should be unique per provider
		"ContentSize",    //Size in bytes
		"ContentType",    //Mime type like image/ext
		"Description",    //Similar to the subject
		"EncodingFormat", //File extension
		"Width",          //Image width
		"Height",         //Image extension
		"Keywords",       //Keywords
		"Name",           //Display name  or Title
	}
}
