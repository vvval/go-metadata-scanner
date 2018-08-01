package pool

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
	"time"
)

const NumHander = 10

var (
	wg     sync.WaitGroup
	jobs   chan string
	status chan string
	errors chan error
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "pool",
		Run: poolHandler,
	})

	jobs = make(chan string)
	status = make(chan string)
	errors = make(chan error)

	// worker pool
	for i := 0; i < NumHander; i++ {
		go func(input <-chan string, status chan string, errors chan error) {
			for {
				select {
				case i, ok := <-input:
					if !ok {
						return
					}

					res, err := work(i)
					if err != nil {
						errors <- err
					} else {
						status <- res
					}

					wg.Done()
				}
			}
		}(jobs, status, errors)
	}
}

func poolHandler(cmd *cobra.Command, args []string) {
	tick := time.Tick(500 * time.Millisecond)
	go func() {
		var okInc, errInc int
		for {
			select {
			case _, ok := <-status:
				if !ok {
					return
				}
				okInc++
				//fmt.Println(s)
			case err := <-errors:
				errInc++
				fmt.Println("ERROR ", err)
			case <-tick:
				fmt.Printf("Ok %d, err %d\n", okInc, errInc)
			}
		}
	}()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		jobs <- strconv.Itoa(i)
	}

	wg.Wait()
	close(jobs)
	close(status)
	close(errors)

	fmt.Println("all done")
}

func work(in string) (res string, err error) {
	time.Sleep(time.Millisecond * 100)

	if in == "500" {
		return "", fmt.Errorf("some error")
	}

	return in, nil
}
