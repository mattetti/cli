package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

var Demo = func(c *cli.Context) {
	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		config, err = configuration.Demo()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	assnmts, err := api.FetchAssignments(config, api.FetchEndpoints["demo"])
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, a := range assnmts {
		if err := a.Save(config.ExercismDirectory); err != nil {
			fmt.Println(err)
		}
	}

	msg := "\nThe demo exercises have been written to %s, in subdirectories by language.\n\nTo try an exercise, change directory to a language/exercise, read the README and run the tests.\n\n"
	fmt.Printf(msg, config.ExercismDirectory)
}
