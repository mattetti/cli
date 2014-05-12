package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

var Restore = func(c *cli.Context) {
	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		fmt.Println("Are you sure you are logged in? Please login again.")
		return
	}

	assnmts, err := api.FetchAssignments(config, api.FetchEndpoints["restore"])
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, a := range assnmts {
		if err := a.Save(config.ExercismDirectory); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Exercises written to %s\n", config.ExercismDirectory)
}
