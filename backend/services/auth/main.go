package main

import (
	"log"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/cmd"
	"github.com/ONLYOFFICE/zoom-onlyoffice/services/auth/config"
)

func main() {
	if err := cmd.Run(config.BuildConfig()); err != nil {
		log.Fatalln(err)
	}
}
