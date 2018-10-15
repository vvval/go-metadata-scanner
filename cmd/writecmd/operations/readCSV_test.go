package operations

import (
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/configuration"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/vars/metadata"
	"testing"
)

func TestReadCSV(t *testing.T) {
	filename := "./fixtures/file1.csv"
	file := util.MustOpenReadonlyFile(filename)
	defer file.Close()

	dict := configuration.Load(config.DictConfig{}, "./../../../dict.yaml").(config.DictConfig)

	ReadCSV(util.GetCSVReader(file, ','), dict, func(filename string, payload metadata.Payload) {
		//wg.Add(1)
		//jobs <- writecmd.NewJob(filename, payload)
	})
}
