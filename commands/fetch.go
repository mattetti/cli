package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

var Fetch = func(c *cli.Context) {
	argCount := len(c.Args())
	if argCount < 0 || argCount > 2 {
		fmt.Println("Usage: exercism fetch\n   or: exercism fetch LANGUAGE\n   or: exercism fetch LANGUAGE EXERCISE")
		return
	}

	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		if argCount == 0 || argCount == 1 {
			fmt.Println("Are you sure you are logged in? Please login again.")
			return
		} else {
			config, err = configuration.Demo()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	assnmts, err := api.FetchAssignments(config, api.FetchEndpoint(c.Args()))
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(assnmts) == 0 {
		noAssignmentMessage := "No assignments found"
		if argCount == 2 {
			fmt.Printf("%s for %s - %s\n", noAssignmentMessage, c.Args()[0], c.Args()[1])
		} else if argCount == 1 {
			fmt.Printf("%s for %s\n", noAssignmentMessage, c.Args()[0])
		} else {
			fmt.Printf("%s\n", noAssignmentMessage)
		}
		return
	}

	for _, a := range assnmts {
		if err := a.Save(config.ExercismDirectory); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("Exercises written to %s\n", config.ExercismDirectory)
}
