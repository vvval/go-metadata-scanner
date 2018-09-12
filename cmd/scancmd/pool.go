package scancmd

import (
	"github.com/vvval/go-metadata-scanner/log"
	"sync"
)

// Create a pool of defined size.
// Read files input chan (files chunk).
// Execute a callback function to read files chunk
// Parse read result and put into output chan (each file separately)
func CreatePool(
	wg *sync.WaitGroup,
	poolSize int,
	chunks <-chan Chunk,
	scanFilesCallback func(files Chunk) ([]byte, error),
	output chan<- FileData,
) {
	for i := 0; i < poolSize; i++ {
		go func(files <-chan Chunk) {
			for {
				select {
				case chunk, ok := <-files:
					if !ok {
						return
					}

					res, err := scanFilesCallback(chunk)
					logWork(res, err)

					if err == nil {
						for _, parsed := range parse(res, chunk) {
							output <- parsed
						}
					}

					wg.Done()
				}
			}
		}(chunks)
	}
}

func logWork(result []byte, err error) {
	if err != nil {
		log.Failure("", err.Error())
	} else if len(result) != 0 {
		//log.Success("Success", string(result))
	}
}
