package cmd

import (
	"fmt"

	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
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

		cmdArgs := []string{}

		for _, k := range appConfig.Fields {
			cmdArgs = append(cmdArgs, fmt.Sprintf("-%s:all", k))
		}
		cmdArgs = append(cmdArgs, "-j", "-G", cmdInput.filename)

		fmt.Println("cmd args: %+v\n", cmdArgs)
		execCmd := exec.Command(appConfig.ExifToolPath, cmdArgs...)
		result, err := execCmd.Output()
		fmt.Println(string(result))
		if err != nil {
			log.Fatal(err)
		} else {
			a := [1]map[string]interface{}{}
			err = json.Unmarshal([]byte(result), &a)
			if err != nil {
				log.Fatal(err)
			}
		}
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
