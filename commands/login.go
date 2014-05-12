package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/configuration"
)

var Login = func(c *cli.Context) {
	config, err := askForConfigInfo()
	if err != nil {
		fmt.Println(err)
		return
	}
	configuration.ToFile(c.GlobalString("config"), config)
	fmt.Printf("Your exercism directory can be found at %s\n", config.ExercismDirectory)
}
