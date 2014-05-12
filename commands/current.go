package commands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

// Current retrieves the current assignments.
var Current = func(c *cli.Context) {
	var language string
	argc := len(c.Args())
	if argc != 0 && argc != 1 {
		fmt.Println("Usage: exercism current\n   or: exercism current LANGUAGE")
		return
	}

	config, err := configuration.FromFile(c.GlobalString("config"))
	if err != nil {
		fmt.Println("Are you sure you are logged in? Please login again.")
		return
	}
	currentAssignments, err := api.FetchAssignments(config, api.FetchEndpoints["current"])
	if err != nil {
		fmt.Println(err)
		return
	}

	if argc == 1 {
		language = c.Args()[0]
		fmt.Println("Current Assignments for", strings.Title(language))
	} else {
		fmt.Println("Current Assignments")
	}

	for _, a := range currentAssignments {
		if argc == 1 {
			if strings.ToLower(language) == strings.ToLower(a.Track) {
				fmt.Printf("%v: %v\n", strings.Title(a.Track), a.Slug)
			}
		} else {
			fmt.Printf("%v: %v\n", strings.Title(a.Track), a.Slug)
		}
	}
}
