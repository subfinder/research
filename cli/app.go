package cli

import "github.com/urfave/cli"

const (
	Name    = "SubFinder"
	Version = "v2.0.0"
)

func NewApplication() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Commands = []cli.Command{}
	return app
}
