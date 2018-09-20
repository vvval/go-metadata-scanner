package operations

import (
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/etool"
	"github.com/vvval/go-metadata-scanner/vars"
	"sync"
)

func ScanFiles(fileChunks []vars.Chunk, poolSize int) []vars.File {
	var done = make(chan struct{})

	var wg sync.WaitGroup
	var chunks = make(chan vars.Chunk)
	var scannedFiles = make(chan vars.File)
	scancmd.CreatePool(
		&wg,
		poolSize,
		chunks,
		func(files vars.Chunk) ([]byte, error) {
			return etool.Read(files, config.App.Fields())
		},
		func(data []byte) {
			for _, parsed := range etool.Parse(data) {
				scannedFiles <- parsed
			}
		},
	)

	for _, chunk := range fileChunks {
		wg.Add(1)
		chunks <- chunk
	}

	go func() {
		wg.Wait()
		close(scannedFiles)
		close(chunks)
		done <- struct{}{}
	}()

	var files []vars.File

	for file := range scannedFiles {
		files = append(files, file)
	}

	<-done
	close(done)

	return files
}
