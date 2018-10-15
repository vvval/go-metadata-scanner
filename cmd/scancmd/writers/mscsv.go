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
	dict   config.DictConfig
	config config.MSCSVConfig
	BaseWriter
	csv *csv.Writer
}

// Headers to be like: Filename, XMP, IPTC, etc...
func (w *MSCSVWriter) Write(file *vars.File) error {
	record, err := packMSCSVLine(file, w.dict, w.config)
	if err != nil {
		return err
	}

	return w.csv.Write(record)
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

func headers() []string {
	return []string{
		"Provider",       //eLearningBrothers
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

func packMSCSVLine(f *vars.File, dict config.DictConfig, config config.MSCSVConfig) ([]string, error) {
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
		config.Provider(),
		f.RelPath(),
		fmt.Sprintf("%d", stat.Size()),
		mime.TypeByExtension(filepath.Ext(file.Name())),
		description(f, dict),
		util.Extension(f.Filename()),
		fmt.Sprintf("%d", img.Width),
		fmt.Sprintf("%d", img.Height),
		keywords(f, dict),
		title(f, dict),
	}

	return record, nil
}

func description(file *vars.File, dict config.DictConfig) string {
	return findTagValue(file, "description", dict).(string)
}

func keywords(file *vars.File, dict config.DictConfig) string {
	keywords := findTagValue(file, "keywords", dict)

	if str, ok := keywords.(string); ok {
		return str
	}

	if isSlice(keywords) {
		return slice2string(keywords)
	}

	return fmt.Sprintf("%v", keywords)
}

func isSlice(i interface{}) bool {
	value := reflect.ValueOf(i)

	return value.Kind() == reflect.Slice
}

func slice2string(i interface{}) string {
	v := reflect.ValueOf(i)
	output := make([]string, v.Len())

	for i := 0; i < v.Len(); i++ {
		output[i] = fmt.Sprintf("%s", v.Index(i))
	}

	return strings.Join(output, ", ")
}

func title(file *vars.File, dict config.DictConfig) string {
	return findTagValue(file, "title", dict).(string)
}

func findTagValue(file *vars.File, field string, dict config.DictConfig) interface{} {
	if tag, found := dict.Find(field); found {
		for _, name := range tag.Map() {
			if value, ok := file.Tags().Tag(name); ok && value != nil && value != "" {
				return value
			}
		}
	}

	return ""
}
