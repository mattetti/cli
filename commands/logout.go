package commands

import "github.com/codegangsta/cli"

var Logout = func(c *cli.Context) {
	logout(c.GlobalString("config"))
}
