package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "github.com/urfave/cli"
)

// Run a command, and fail if an error is encountered
func runCommandOrFail(cmd string, args []string, printOutput bool) {
  output, err := exec.Command(cmd, args...).CombinedOutput()
  if len(output) > 0 && printOutput {
    fmt.Println(string(output))
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

// Clone the Kafka Streams skeleton project
func cloneSkeletonRepo(intoDir string) {
  cmd := "git"
  repoUrl := ""
  args := []string{"clone", repoUrl, intoDir}
  runCommandOrFail(cmd, args, false)
  os.Chdir(intoDir)
}

// Initialize a new git repository
func initializeNewGitRepo(dir string) {
  cloneSkeletonRepo(dir)

  // Remove the old .git files
  runCommandOrFail("rm", []string{"-rf", ".git"}, false)

  // Initialize the new repo
  runCommandOrFail("git", []string{"init"}, false)
}

// Walk the file path and removes placeholders from files and directories
func removePlaceholdersFromFiles(projectName string) {
  dir, err := os.Getwd()
  if (err != nil) {
    fmt.Println(err)
    os.Exit(1)
  }
  // rename directories first
  filepath.Walk(dir, func(path string, fileInfo os.FileInfo, err error) (e error) {
    return renameFilesAndDirectories(path, projectName, fileInfo, err)
  })

  filepath.Walk(dir, func(path string, fileInfo os.FileInfo, err error) (e error) {
    return removePlaceholders(path, projectName, fileInfo, err)
  })
}

// Rename files and directories containing the placeholder value
func renameFilesAndDirectories(path string, projectName string, fileInfo os.FileInfo, err error) (e error) {
  dir := filepath.Dir(path)
  // Rename files and directories
  if strings.HasPrefix(fileInfo.Name(), "myproject") {
    base := filepath.Base(path)
    newFileName := filepath.Join(dir, strings.Replace(base, "myproject", projectName, 1))
    os.Rename(path, newFileName)
  }
  return
}

// Remove placeholders from the file's content
func removePlaceholders(path string, projectName string, fileInfo os.FileInfo, err error) (e error) {
  dir := filepath.Dir(path)
  // Search and replace the placeholder text in files
  if (!fileInfo.IsDir()) {
    searchFileAndReplace(filepath.Join(dir, fileInfo.Name()), fileInfo.Mode(), "myproject", projectName)
  }
  return
}

// Search a file and replace all instances of a target string with a new string
func searchFileAndReplace(path string, mode os.FileMode, search string, replace string) {
  input, err := ioutil.ReadFile(path)
  if (err != nil) {
    fmt.Println(err)
    os.Exit(1)
  }

  output := bytes.Replace(input, []byte(search), []byte(replace), -1)
  if err = ioutil.WriteFile(path, output, mode); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

// Create the first commit in the new Kafka Streams project
func firstCommit() {
  runCommandOrFail("git", []string{"add", "."}, false)
  runCommandOrFail("git", []string{"commit", "-m", "Project initialized with stream-utils"}, false)
}

// Print the welcome message
func printWelcome(projectName string) {
  bytes, err := ioutil.ReadFile("welcome.txt") // just pass the file name
  if err != nil {
    fmt.Println("Could not print welcome message :(")
  }

  // Print the contents of the welcome file
  welcome := string(bytes)
  fmt.Println(welcome)

}

// Create a new CLI app
func main() {
  app := cli.NewApp()
  app.Name = "stream-utils"
  app.Usage = "Utilities for creating a Kafka Streams skeleton project"
  app.Version = "0.1.0"
  app.Commands = []cli.Command{
    {
      Name:    "create",
      Aliases: []string{"c"},
      Usage:   "Creates a new skeleton project",
      Action:  func(c *cli.Context) error {
        newProjectName := c.Args().Get(0)
        initializeNewGitRepo(newProjectName)
        removePlaceholdersFromFiles(newProjectName)
        firstCommit()
        printWelcome(newProjectName)
        return nil
      },
    },
  }
  app.Run(os.Args)
}
