package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/configuration"
)

var Whoami = func(c *cli.Context) {
	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		fmt.Println("Are you sure you are logged in? Please login again.")
		return
	}

	fmt.Println(config.GithubUsername)
}
