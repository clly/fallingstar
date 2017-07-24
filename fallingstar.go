package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/github"
)

func main() {
	var user string
	if len(os.Args) >= 2 {
		user = os.Args[1]
	}
	client := github.NewClient(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	starredRepos, resp, err := client.Activity.ListStarred(ctx, user, nil)
	if err != nil {
		fmt.Errorf("Failed to pull starred repos: %s", err)
	}
	for i := range starredRepos {
		repo := starredRepos[i].Repository
		fullName := *repo.FullName
		_, err := os.Stat(fullName)
		if err != nil {
			fmt.Printf("Going to clone into %s via %s\n", fullName, *repo.CloneURL)
			clone(fullName, *repo.CloneURL)
			time.Sleep(time.Second * 5)
		} else {
			fmt.Printf("Unable to clone into path %s. It already exists\n", fullName)
			pull(fullName)
			time.Sleep(time.Second * 1)
		}
	}
	fmt.Printf("Github api limit: %v. Github api remaining: %v. Github api reset %v\n", resp.Limit, resp.Remaining, resp.Reset)
}

func pull(path string) {
	command := exec.Command("git", "pull", "--ff-only")
	command.Dir = path
	// TODO: Replicating from git clone. Combine functions
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Pulling updates from upstream\n%s\n", out.String())
}

func clone(path string, url string) {
	command := exec.Command("git", "clone", url, path)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cloning url into path: %s\n%s\n", path, out.String())
}
