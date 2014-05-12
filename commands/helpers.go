package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/exercism/cli/configuration"
)

var (
	testExtensions = map[string]string{
		"ruby":    "_test.rb",
		"js":      ".spec.js",
		"elixir":  "_test.exs",
		"clojure": "_test.clj",
		"python":  "_test.py",
		"go":      "_test.go",
		"haskell": "_test.hs",
	}
)

func logout(path string) {
	os.Remove(path)
}

func absolutePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

func askForConfigInfo() (c configuration.Config, err error) {
	var un, key, dir string
	delim := "\r\n"

	bio := bufio.NewReader(os.Stdin)

	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	fmt.Print("Your GitHub username: ")
	un, err = bio.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Print("Your exercism.io API key: ")
	key, err = bio.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Println("What is your exercism exercises project path?")
	fmt.Printf("Press Enter to select the default (%s):\n", currentDir)
	fmt.Print("> ")
	dir, err = bio.ReadString('\n')
	if err != nil {
		return
	}

	key = strings.TrimRight(key, delim)
	un = strings.TrimRight(un, delim)
	dir = strings.TrimRight(dir, delim)

	if dir == "" {
		dir = currentDir
	}

	dir = configuration.ReplaceTilde(dir)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = fmt.Errorf("Error making directory %v: [%v]", dir, err)
		return
	}

	dir, err = absolutePath(dir)
	if err != nil {
		return
	}

	c = configuration.Config{GithubUsername: un, ApiKey: key, ExercismDirectory: dir, Hostname: "http://exercism.io"}
	return
}

// isTest checks if a filename refers to a test file.
func isTest(filename string) bool {
	for _, ext := range testExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
