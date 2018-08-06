package scancmd

import (
	"github.com/vvval/go-metadata-scanner/log"
	"sync"
)

func CreatePool(wg *sync.WaitGroup, poolSize int, files <-chan []string, callback func(files []string) ([]byte, error)) {
	for i := 0; i < poolSize; i++ {
		go func(files <-chan []string) {
			for {
				select {
				case chunk, ok := <-files:
					if !ok {
						return
					}

					res, err := callback(chunk)
					logWork(res, err)

					wg.Done()
				}
			}
		}(files)
	}
}

func logWork(result []byte, err error) {
	if err != nil {
		log.Failure("", err.Error())
	} else if len(result) != 0 {
		log.Success("Success", string(result))
	}
}
