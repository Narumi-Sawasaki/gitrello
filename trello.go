package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

var getCurrentBranchName = func() (string, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}
	ref, err := r.Head()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return ref.Name().Short(), nil
}

func getCurrentTaskName() (string, error) {
	branchName, err := getCurrentBranchName()
	if err != nil {
		return "", err
	}

	prefix := "task/"
	if len(branchName) < len(prefix) || branchName[:len(prefix)] != prefix {
		return "", errors.New("ブランチ名が不正です")
	}

	return branchName[len(prefix):], nil
}
