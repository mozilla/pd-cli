package main

import (
	"os"

	"github.com/mozilla/pd-cli/command/repo"
	"gopkg.in/urfave/cli.v1"
)

func main() {

	app := cli.NewApp()
	app.Commands = []cli.Command{
		repo.NewCommand(),
	}
	app.Run(os.Args)
}
