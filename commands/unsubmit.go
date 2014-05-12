package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

var Unsubmit = func(c *cli.Context) {
	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		fmt.Println("Are you sure you are logged in? Please login again.")
		return
	}

	response, err := api.UnsubmitAssignment(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response != "" {
		return
	}

	fmt.Println("The last submission was successfully deleted.")
}
