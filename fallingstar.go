package main

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"log"
	"time"
)

func main() {
	client := github.NewClient(nil)
	starredRepos, resp, err := client.Activity.ListStarred("clly", nil)
	if err != nil {
		panic(err)
	}
	for i := range(starredRepos) {
		repo := starredRepos[i].Repository
		fullName := *repo.FullName
		fmt.Printf("Going to clone into %s via %s\n", fullName, *repo.CloneURL)
		_, err := os.Stat(fullName)
		if err != nil {
			clone(*repo.FullName, *repo.CloneURL)
			time.Sleep(time.Second * 5)
		} else {
			fmt.Printf("Unable to clone into path %s. It already exists\n", fullName)
		}
	}
	fmt.Println(resp)
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