package main

import (
	"github.com/google/go-github/github"
	"fmt"
	"os/exec"
)

func main() {
	client := github.NewClient(nil)
	//user, resp, err:= client.Users.Get("clly")
	starredRepos, resp, err := client.Activity.ListStarred("clly", nil)
	if err != nil {
		panic(err)
	}
	for i := range(starredRepos) {
		repo := starredRepos[i].Repository

		fmt.Println(*repo.Name, *repo.CloneURL)
	}
	fmt.Println(resp)
}

func clone(name string, url string) {
	command := fmt.Sprintf("git clone %s %s", name, url)
	exec.Cmd(command)
}