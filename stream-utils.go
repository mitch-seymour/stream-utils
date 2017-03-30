package main

import (
  "fmt"
  "os"
  "os/exec"

  "github.com/urfave/cli"
)

func runCommandOrFail(cmd string, args []string) {
  err := exec.Command(cmd, args...).Run()
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func cloneSkeletonRepo(intoDir string) {
  fmt.Println("Cloning skeleton project")
  cmd := "git"
  args := []string{"clone", "git@github.com:mitch-seymour/retro-futurism", intoDir}
  runCommandOrFail(cmd, args)
  os.Chdir(intoDir)
}

func initializeNewGitRepo(dir string) {
  fmt.Println("Initializing git repo")
  // first, remove the old git repo
  runCommandOrFail("rm", []string{"-rf", ".git"})
  // now, initialize the new repo
  runCommandOrFail("git", []string{"init"})
}

func main() {
  app := cli.NewApp()
  app.Name = "stream-utils"
  app.Usage = "Utilities for creating a Kafka Streams skeleton project"
  app.Version = "0.1.0"
  app.Commands = []cli.Command{
    {
      Name:    "create",
      Aliases: []string{"h"},
      Usage:   "Creates a new skeleton project",
      Action:  func(c *cli.Context) error {
        newProjectName := c.Args().Get(0)
        fmt.Printf("Creating new Kafka Streams app: %s\n", newProjectName)
        cloneSkeletonRepo(newProjectName)
        initializeNewGitRepo(newProjectName)
        return nil
      },
    },
  }
  app.Run(os.Args)
}
