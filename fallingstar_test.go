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
