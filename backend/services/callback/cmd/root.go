package cmd

import (
	"os"

	"github.com/ONLYOFFICE/zoom-onlyoffice/services/callback/config"
	"github.com/urfave/cli/v2"
)

func GetCommands(config *config.Config) cli.Commands {
	return []*cli.Command{
		Server(config),
	}
}

func Run(config *config.Config) error {
	app := &cli.App{
		Name:        "onlyoffice:callback",
		Description: "Description",
		Authors: []*cli.Author{
			{
				Name:  "Ascensio Systems SIA",
				Email: "support@onlyoffice.com",
			},
		},
		HideVersion: true,
		Commands:    GetCommands(config),
	}

	return app.Run(os.Args)
}
