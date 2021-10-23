package main

import (
	"log"
	"os"

	//X- range .Modules
	"package_name/internal/module_name" /*X end X*/
)

func run(_ []string) error {
	/*X range .Modules X*/
	module_name.SaySomething()
	/*X end X*/
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
