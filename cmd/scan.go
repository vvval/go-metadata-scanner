package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/cmd/scancmd"
	"github.com/vvval/go-metadata-scanner/config"
	"github.com/vvval/go-metadata-scanner/util"
	"math"
	"sync"
)

func init() {
	// cmd represents the scan command
	var cmd = &cobra.Command{
		Use:   "scan",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: scanHandler,
	}

	rootCmd.AddCommand(cmd)
	scancmd.FillInput(cmd, &scancmd.Input)

	scancmd.Files = make(chan []string)
	scancmd.CreatePool(&wg, scancmd.PoolSize, scancmd.Files, exec)
}

func exec(names []string) ([]byte, error) {
	var args = []string{"-j", "-G"}

	for _, field := range config.Get().Fields() {
		args = append(args, fmt.Sprintf("-%s:all", field))
	}

	for _, name := range names {
		args = append(args, name)
	}

	return util.RunCommand(config.Get().ToolPath(), args...)
}

func scanHandler(cmd *cobra.Command, args []string) {
	ReadJSON(scancmd.Input.Directory(), config.Get().Extensions(), scancmd.PoolSize, &wg, scancmd.Files)

	wg.Wait()
	close(scancmd.Files)
	//var m = map[string]interface{}{}
	//var a = []interface{}{}
	//execCmd := exec.Command(config.Get().ToolPath(), "-j", "-IPTC:Keywords", "-XMP:Subject", "-XMP:Marked", "-IPTC:Headline", "-XMP:Headline", "keywords/test.jpg", "keywords/001.png")
	//result, err := execCmd.Output()
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	a := make([]map[string]interface{}, 2)
	//	err = json.Unmarshal([]byte(result), &a)
	//	fmt.Printf("%s\n", result)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//fmt.Printf("%+v\n", a)
	//	//for _, m := range a {
	//	//fmt.Printf("%s\n", json.MarshalIndent(m,"","  "))
	//	//for k, v := range m {
	//	//fmt.Printf("%s:%v\n", k, v)
	//	////or via type assertion
	//	//switch reflect.TypeOf(v).Kind() {
	//	//case reflect.Array:
	//	//	fallthrough
	//	//case reflect.Slice:
	//	//	s := reflect.ValueOf(v)
	//	//
	//	//	for i := 0; i < s.Len(); i++ {
	//	//		fmt.Printf("%d: %v\n\n\n", i, s.Index(i))
	//	//	}
	//	//}
	//	////for vk, vv := range v {
	//	////	fmt.Printf("%+v:%+v\n", vk, vv)
	//	////}
	//	//}
	//	//}
	//}
	//
	////a = append(a, m)
	//
	//return

	//writeInput := writeCommand.Flags()
	//cmdArgs := []string{}
	//
	//for _, k := range config.Get().Fields {
	//	cmdArgs = append(cmdArgs, fmt.Sprintf("-%s:all", k))
	//}
	//cmdArgs = append(cmdArgs, "-j", "-G", writeInput.filename())
	//
	//fmt.Println("cmd args: %+v\n", cmdArgs)
	//execCmd := exec.Command(config.Get().ToolPath, cmdArgs...)
	//result, err := execCmd.Output()
	//fmt.Println(string(result))
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	a := [1]map[string]interface{}{}
	//	err = json.Unmarshal([]byte(result), &a)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
}

func ReadJSON(directory string, extensions []string, poolSize int, wg *sync.WaitGroup, out chan<- []string) {
	files = scanDir(directory, extensions)
	var chunkSize = int(math.Ceil(float64(len(files) / poolSize)))

	for i := 0; i < len(files); i += chunkSize {
		end := i + chunkSize
		if end > len(files) {
			end = len(files)
		}

		wg.Add(1)

		chunk := files[i:end]
		if len(chunk) > 0 {
			out <- chunk
		}
	}
}
