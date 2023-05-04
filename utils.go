package main

import "github.com/go-git/go-git/v5"

func getRepoRoot() (string, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	return w.Filesystem.Root(), nil
}
