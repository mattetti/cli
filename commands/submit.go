package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/exercism/cli/api"
	"github.com/exercism/cli/configuration"
)

var Submit = func(c *cli.Context) {
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

}
