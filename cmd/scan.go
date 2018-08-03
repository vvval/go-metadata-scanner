package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vvval/go-metadata-scanner/config"
	"log"
	"os/exec"
	"reflect"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//var m = map[string]interface{}{}
		//var a = []interface{}{}
		execCmd := exec.Command(config.Get().ToolPath(), "-j", "-Keywords", "keywords/test.jpg")
		result, err := execCmd.Output()
		if err != nil {
			log.Fatal(err)
		} else {
			a := [1]map[string]interface{}{}
			err = json.Unmarshal([]byte(result), &a)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%+v\n", a)
			for k, v := range a[0] {
				fmt.Printf("%s:%v\n", k, v)
				//or via type assertion
				switch reflect.TypeOf(v).Kind() {
				case reflect.Array:
					fallthrough
				case reflect.Slice:
					s := reflect.ValueOf(v)

					for i := 0; i < s.Len(); i++ {
						fmt.Printf("%d: %v\n", i, s.Index(i))
					}
				}
				//for vk, vv := range v {
				//	fmt.Printf("%+v:%+v\n", vk, vv)
				//}
			}
		}

		//a = append(a, m)

		return
		//input := writer.Input()
		//cmdArgs := []string{}
		//
		//for _, k := range config.Get().Fields {
		//	cmdArgs = append(cmdArgs, fmt.Sprintf("-%s:all", k))
		//}
		//cmdArgs = append(cmdArgs, "-j", "-G", input.filename())
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
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	//fmt.Printf("%+v\n", appConfig)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//scanCmd.Flags().StringP("dir","d", "", "Directory to scan")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
