package scancmd

import (
	"github.com/vvval/go-metadata-scanner/util/log"
	"github.com/vvval/go-metadata-scanner/vars"
	"sync"
)

// Create a pool of defined size.
// Read files input chan (files chunk).
// Execute a callback function to read files chunk
// Parse read result and put into output chan (each file separately)
func CreatePool(
	wg *sync.WaitGroup,
	poolSize int,
	chunks <-chan vars.Chunk,
	scanFilesCallback func(files vars.Chunk) ([]byte, error),
	output chan<- vars.File,
) {
	for i := 0; i < poolSize; i++ {
		go func(files <-chan vars.Chunk) {
			for {
				select {
				case chunk, ok := <-files:
					if !ok {
						return
					}

					res, err := scanFilesCallback(chunk)
					if err != nil {
						log.Failure("Scan error", err.Error())
					} else {
						for _, parsed := range parse(res) {
							output <- parsed
						}
					}

					wg.Done()
				}
			}
		}(chunks)
	}
}
