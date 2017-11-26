package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-github/github"
)

const binGit = "/bin/git"
const usrBinGit = "/usr/bin/git"
const usrLocalGit = "/usr/local/bin/git"

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
		fmt.Fprintf(os.Stderr, "Failed to pull starred repos: %s\n", err)
	}
	for i := range starredRepos {
		repo := starredRepos[i].Repository
		fullName := *repo.FullName
		_, err := os.Stat(fullName)
		if err != nil {
			fmt.Printf("Going to clone into %s via %s\n", fullName, *repo.CloneURL)
			err = execGitCmd(fullName, *repo.CloneURL)
			if err != nil {
				fmt.Println("Failed to clone ", fullName, err)
			}
			time.Sleep(time.Second * 5)
		} else {
			fmt.Printf("Unable to clone into path %s. It already exists\n", fullName)
			err = execGitCmd(fullName, "")
			if err != nil {
				fmt.Println("Failed to pull git repo ", err)
			}
			time.Sleep(time.Second * 1)
		}
	}
	fmt.Printf("Github api limit: %v. Github api remaining: %v. Github api reset %v\n", resp.Limit, resp.Remaining, resp.Reset)
}

// buildGitCmd creates the git clone or git pull command. If url is empty
// string it creates a git pull instead of git clone
func gitCmd(path string, url string) (*exec.Cmd, error) {
	gitPath := findGit()
	if gitPath == "" {
		return nil, fmt.Errorf("Unable to find git on PATH or %s %s %s", binGit, usrBinGit, usrLocalGit)
	}
	var pwd string
	gitCmd := make([]string, 0, 4)
	if url != "" {
		gitCmd = []string{gitPath, "clone", url, path}
	} else {
		gitCmd = []string{gitPath, "pull", "--ff-only"}
		pwd = path
	}

	command := exec.Command(gitPath, gitCmd[1:]...) // #nosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stdout
	// Should we make this the fully qualified path?
	command.Dir = pwd
	return command, nil
}

// execGitCmd is a wrapper around cmd.Start from gitCmd. It makes testing easier
func execGitCmd(path, url string) error {
	cmd, err := gitCmd(path, url)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func findGit() string {
	if exists(binGit) {
		return binGit
	} else if exists(usrBinGit) {
		return usrBinGit
	} else if exists(usrLocalGit) {
		return usrLocalGit
	}
	gitPath, err := exec.LookPath("git")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to find git on path. Looking at some other places")
		return ""
	}
	return gitPath
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return true
}
