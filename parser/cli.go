package parser

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var supportedParamFile = map[string]struct{}{".json": {}, ".yaml": {}, ".yml": {}}

func loadParam(paramFile string) map[string]interface{} {
	var res map[string]interface{}
	if paramFile == "" {
		return res
	}
	var ext = filepath.Ext(paramFile)
	if _, ok := supportedParamFile[ext]; !ok {
		return res
	}

	data, err := ioutil.ReadFile(paramFile)
	if err != nil {
		log.Println("Reading param file error", err)
		return res
	}
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &res); err != nil {
			log.Println("Reading param file error", err)
		}
	case ".yaml":
		fallthrough
	case ".yml":
		if err := yaml.Unmarshal(data, &res); err != nil {
			log.Println("Reading param file error", err)
		}
	}
	return res
}
func generate(source, destination, paramFile string) {
	var params = loadParam(paramFile)
	var parser = New(
		WithFS(os.DirFS(".")),
		WithParams(params),
	)
	if err := parser.Parse(source, destination); err != nil {
		log.Fatal(err)
	}
}

func Command() *cobra.Command {
	var destination = "./parsed"
	var paramFile string
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Convert your prototype into reproducible golang templates",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Printf("Dunno")
				return
			}
			var source = args[0]
			if len(args) > 1 {
				destination = args[1]
			}
			generate(source, destination, paramFile)
		},
	}
	cmd.Flags().StringVarP(&paramFile, "param", "p", "", "Parameter file for rebuild the dismantled project")
	return cmd
}
