package main

import (
	"log"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/cmd"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/config"
)

func main() {
	if err := cmd.Run(config.BuildConfig()); err != nil {
		log.Fatalln(err)
	}
}
