package main

import (
	"log"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/cmd"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/gateway/config"
)

func main() {
	if err := cmd.Run(config.BuildConfig()); err != nil {
		log.Fatalln(err)
	}
}
