package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/assignments"
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
			Action: func(c *cli.Context) {
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
			},
		},
		{
			Name:      "demo",
			ShortName: "d",
			Usage:     "Fetch first assignment for each language from exercism.io",
			Action: func(c *cli.Context) {
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
					err := assignments.Save(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				msg := "\nThe demo exercises have been written to %s, in subdirectories by language.\n\nTo try an exercise, change directory to a language/exercise, read the README and run the tests.\n\n"
				fmt.Printf(msg, config.ExercismDirectory)
			},
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Usage:     "Fetch assignments from exercism.io",
			Action: func(c *cli.Context) {
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
					err := assignments.Save(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", config.ExercismDirectory)
			},
		},
		{
			Name:      "login",
			ShortName: "l",
			Usage:     "Save exercism.io api credentials",
			Action: func(c *cli.Context) {
				config, err := askForConfigInfo()
				if err != nil {
					fmt.Println(err)
					return
				}
				configuration.ToFile(c.GlobalString("config"), config)
				fmt.Printf("Your exercism directory can be found at %s\n", config.ExercismDirectory)
			},
		},
		{
			Name:      "logout",
			ShortName: "o",
			Usage:     "Clear exercism.io api credentials",
			Action: func(c *cli.Context) {
				logout(c.GlobalString("config"))
			},
		},
		{
			Name:      "restore",
			ShortName: "r",
			Usage:     "Restore completed and current assnmts from exercism.io",
			Description: "Restore will pull the latest revisions of exercises that have already been " +
				"submitted. It will *not* overwrite existing files.  If you have made changes " +
				"to a file and have not submitted it, and you're trying to restore the last " +
				"submitted version, first move that file out of the way, then call restore.",
			Action: func(c *cli.Context) {
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
					err := assignments.Save(config.ExercismDirectory, a)
					if err != nil {
						fmt.Println(err)
					}
				}

				fmt.Printf("Exercises written to %s\n", config.ExercismDirectory)
			},
		},
		{
			Name:      "submit",
			ShortName: "s",
			Usage:     "Submit code to exercism.io on your current assignment",
			Action: func(c *cli.Context) {
				config, err := configuration.FromFile(c.GlobalString("config"))
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}

				if len(c.Args()) == 0 {
					fmt.Println("Please enter a file name")
					return
				}

				filename := c.Args()[0]

				// Make filename relative to config.ExercismDirectory.
				absPath, err := absolutePath(filename)
				if err != nil {
					fmt.Printf("Couldn't find %v: %v\n", filename, err)
					return
				}
				exDir := config.ExercismDirectory + string(filepath.Separator)
				if !strings.HasPrefix(absPath, exDir) {
					fmt.Printf("%v is not under your exercism project path (%v)\n", absPath, exDir)
					return
				}
				filename = absPath[len(exDir):]

				if IsTest(filename) {
					fmt.Println("It looks like this is a test, please enter an example file name.")
					return
				}

				code, err := ioutil.ReadFile(absPath)
				if err != nil {
					fmt.Printf("Error reading %v: %v\n", absPath, err)
					return
				}

				response, err := api.SubmitAssignment(config, filename, code)
				if err != nil {
					fmt.Printf("There was an issue with your submission: %v\n", err)
					return
				}

				fmt.Printf("For feedback on your submission visit %s%s%s\n",
					config.Hostname, "/submissions/", response.Id)

			},
		},
		{
			Name:      "unsubmit",
			ShortName: "u",
			Usage:     "Delete the last submission",
			Action: func(c *cli.Context) {
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
			},
		},
		{
			Name:      "whoami",
			ShortName: "w",
			Usage:     "Get the github username that you are logged in as",
			Action: func(c *cli.Context) {
				config, err := configuration.FromFile(c.GlobalString("config"))
				if err != nil {
					fmt.Println("Are you sure you are logged in? Please login again.")
					return
				}

				fmt.Println(config.GithubUsername)
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
		os.Exit(1)
	}
}
