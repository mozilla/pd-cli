package repo

import (
	"fmt"

	"github.com/urfave/cli"
)

func initRepo(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, "Initializing Repository %s/%s\n", owner, repo)
	return nil
}
