package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/commands"
	"github.com/exercism/cli/configuration"
)

func main() {
	app := cli.NewApp()
	app.Name = "exercism"
	app.Usage = "A command line tool to interact with http://exercism.io"
	app.Version = api.VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", configuration.Filename(configuration.HomeDir()), "path to config file"},
	}
	app.Commands = []cli.Command{
		{
			Name:      "current",
			ShortName: "c",
			Usage:     "Show the current assignments",
			Action:    commands.Current,
		},
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     "Fetch first assignment for each language from exercism.io",
			Action:    commands.Demo,
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     "Fetch assignments from exercism.io",
			Action:    commands.Fetch,
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Save exercism.io api credentials",
			Action:    commands.Login,
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     "Clear exercism.io api credentials",
			Action:    commands.Logout,
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore completed and current assnmts from exercism.io",
			Description: "Restore will pull the latest revisions of exercises that have already been " +
				"submitted. It will *not* overwrite existing files.  If you have made changes " +
				"to a file and have not submitted it, and you're trying to restore the last " +
				"submitted version, first move that file out of the way, then call restore.",
			Action: commands.Restore,
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     "Submit code to exercism.io on your current assignment",
			Action:    commands.Submit,
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     "Delete the last submission",
			Action:    commands.Unsubmit,
		},
		{
			Name:      "whoami",
			ShortName: "w",
			Usage:     "Get the github username that you are logged in as",
			Action:    commands.Whoami,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}
