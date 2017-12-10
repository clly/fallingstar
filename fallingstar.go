package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"gitlab.com/clly/toolbox"

	"github.com/google/go-github/github"
)

const binGit = "/bin/git"
const usrBinGit = "/usr/bin/git"
const usrLocalGit = "/usr/local/bin/git"
const perPage = 100

// Star holds a bunch of
type Star struct {
	Page      int
	NextPage  bool
	User      string
	Limit     int
	Remaining int
	LastPage  int
	Reset     github.Timestamp
}

func main() {
	var user string
	if len(os.Args) >= 2 {
		user = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "Must supply github username")
		os.Exit(1)
	}

	s := &Star{
		User:     user,
		Page:     1,
		NextPage: true,
	}

	for s.NextPage {
		r, err := s.getRepos(s.Page, perPage)
		if err != nil {
			toolbox.Oopse(err)
		}
		fmt.Fprintf(os.Stderr, s.status())
		loopStarred(r)
		fmt.Printf("\nGithub api limit: %v. Github api remaining: %v. Github api reset %v\n",
			s.Limit, s.Remaining, s.Reset)
	}
}

func (s *Star) getRepos(page, perpage int) ([]*github.StarredRepository, error) {
	client := github.NewClient(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	starredRepos, resp, err := client.Activity.ListStarred(ctx, s.User, &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			PerPage: perpage,
			Page:    page,
		},
	})
	if resp.NextPage == 0 {
		s.NextPage = false
	} else {
		s.NextPage = true
	}
	s.Page = resp.NextPage
	s.Reset = resp.Reset
	s.Limit = resp.Limit
	s.Remaining = resp.Remaining
	s.LastPage = resp.LastPage
	return starredRepos, err
}

func (s *Star) status() string {
	if s.Page == 0 {
		return "On last page"
	}
	return fmt.Sprintf("On page %d of %d\n", s.Page-1, s.LastPage)
}

func loopStarred(starredRepos []*github.StarredRepository) {
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
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
