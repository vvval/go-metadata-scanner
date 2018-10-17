package writecmd

import (
	"fmt"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/util"
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/util/scan"
	"github.com/vvval/go-metadata-scanner/vars"
)

func Work(job *Job, append, originals bool, extensions []string, files *vars.Chunk, filesData *[]vars.File) ([]byte, error) {
	filename, found := scan.Candidates(job.Filename(), files, extensions)
	if !found {
		return []byte{}, nil
	}

	if append {
		if file, found := findScanned(filename, filesData); found {
			job.MergePayload(file.Tags())
		}
	}

	if !job.HasPayload() {
		log.Debug("Skip file", fmt.Sprintf("no payload found for `%s`", filename))

		return []byte{}, nil
	}

	payload := job.Payload()
	result, err := etool.Write(
		filename,
		payload.Tags(),
		payload.UseSeparator(),
		originals,
	)

	return result, err
}

func findScanned(filename string, files *[]vars.File) (vars.File, bool) {
	for _, file := range *files {
		if util.PathsEqual(file.Filename(), filename) {
			return file, true
		}
	}

	return vars.File{}, false
}
