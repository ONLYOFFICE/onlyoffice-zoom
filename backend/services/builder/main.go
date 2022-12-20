package main

import (
	"log"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/cmd"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/builder/config"
)

func main() {
	if err := cmd.Run(config.BuildConfig()); err != nil {
		log.Fatalln(err)
	}
}
