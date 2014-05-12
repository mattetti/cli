package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/exercism/cli/configuration"
	"github.com/stretchr/testify/assert"
)

func assertFileDoesNotExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)

	if err == nil {
		t.Errorf("File [%s] already exist.", filename)
	}
}

func TestLogoutDeletesConfigFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)

	c := configuration.Config{}

	configuration.ToFile(tmpDir, c)

	logout(tmpDir)

	assertFileDoesNotExist(t, configuration.Filename(tmpDir))
}

func TestAskForConfigInfoAllowsSpaces(t *testing.T) {
	oldStdin := os.Stdin
	dirName := "dirname with spaces"
	userName := "TestUsername"
	apiKey := "abc123"

	fakeStdin, err := ioutil.TempFile("", "stdin_mock")
	assert.NoError(t, err)

	fakeStdin.WriteString(fmt.Sprintf("%s\r\n%s\r\n%s\r\n", userName, apiKey, dirName))
	assert.NoError(t, err)

	_, err = fakeStdin.Seek(0, os.SEEK_SET)
	assert.NoError(t, err)

	defer fakeStdin.Close()

	os.Stdin = fakeStdin

	c, err := askForConfigInfo()
	if err != nil {
		t.Errorf("Error asking for configuration info [%v]", err)
	}
	os.Stdin = oldStdin
	absoluteDirName, _ := absolutePath(dirName)
	_, err = os.Stat(absoluteDirName)
	if err != nil {
		t.Errorf("Excercism directory [%s] was not created.", absoluteDirName)
	}
	os.Remove(absoluteDirName)
	os.Remove(fakeStdin.Name())

	assert.Equal(t, c.ExercismDirectory, absoluteDirName)
	assert.Equal(t, c.GithubUsername, userName)
	assert.Equal(t, c.ApiKey, apiKey)
}

var isTestTests = []struct {
	filename string
	expected bool
}{
	{"bob_test.rb", true},
	{"bob.spec.js", true},
	{"bob_test.exs", true},
	{"bob_test.clj", true},
	{"bob_test.py", true},
	{"bob_test.go", true},
	{"bob_test.hs", true},
	{"bob.rb", false},
}

func TestIsTest(t *testing.T) {
	for _, test := range isTestTests {
		result := IsTest(test.filename)
		if test.expected != result {
			t.Errorf("Filename [%s] should be a test file but is not.", test.filename)
		}
	}
}
