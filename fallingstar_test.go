package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitCmd(t *testing.T) {
	cmd, err := gitCmd("foo", "bar")
	require.NoError(t, err, "Error building git command")
	args := cmd.Args
	require.Equal(t, "git", filepath.Base(args[0]), "Git is not part of args")
	require.Equal(t, []string{"clone", "bar", "foo"}, args[1:], "Clone args don't match")

	cmd, err = gitCmd("foo", "")
	require.NoError(t, err, "Error building git command")
	args = cmd.Args
	require.Equal(t, []string{"pull", "--ff-only"}, args[1:], "Pull args don't match")
	require.Equal(t, "foo", cmd.Dir, "PWD does not match")
}

func TestGetRepos(t *testing.T) {
	s := Star{
		User: "clly",
		Page: 1,
	}

	repos, err := s.getRepos(s.Page, 1)
	require.NoError(t, err, "Failed to request repositories from github")
	require.Len(t, repos, 1, "Failed to get only 1 starred repo")
	require.True(t, s.NextPage, "There is no next page even though we expect to have more than 1. "+
		"This could be because the user used for the test does not have more than 1 starred repository")
}
