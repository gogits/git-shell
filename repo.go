// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"container/list"
	"errors"
	"os"
	"path"
	"path/filepath"
)

// Repository represents a Git repository.
type Repository struct {
	Path string

	commitCache map[sha1]*Commit
}

const _PRETTY_LOG_FORMAT = `--pretty=format:%H`

func (repo *Repository) parsePrettyFormatLogToList(logs []byte) (*list.List, error) {
	l := list.New()
	if len(logs) == 0 {
		return l, nil
	}

	parts := bytes.Split(logs, []byte{'\n'})

	for _, commitId := range parts {
		commit, err := repo.GetCommit(string(commitId))
		if err != nil {
			return nil, err
		}
		l.PushBack(commit)
	}

	return l, nil
}

// InitRepository initializes a new Git repository.
func InitRepository(repoPath string, bare bool) error {
	os.MkdirAll(repoPath, os.ModePerm)

	cmd := NewCommand("init")
	if bare {
		cmd.AddArguments("--bare")
	}
	_, err := cmd.RunInDir(repoPath)
	return err
}

// OpenRepository opens the repository at the given path.
func OpenRepository(repoPath string) (*Repository, error) {
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	} else if !isDir(repoPath) {
		return nil, errors.New("no such file or directory")
	}

	return &Repository{Path: repoPath}, nil
}

// Clone clones original repository to target path.
func Clone(from, to string) error {
	toDir := path.Dir(to)
	os.MkdirAll(toDir, os.ModePerm)

	_, err := NewCommand("clone", from, to).Run()
	return err
}

// Pull pulls changes from remotes.
func Pull(repoPath string, all bool) error {
	cmd := NewCommand("pull")
	if all {
		cmd.AddArguments("--all")
	}
	_, err := cmd.RunInDir(repoPath)
	return err
}

// Push pushs local commits to given remote branch.
func Push(repoPath, remote, branch string) error {
	_, err := NewCommand("push", remote, branch).RunInDir(repoPath)
	return err
}

// Reset resets HEAD to given revision or head of branch.
func Reset(repoPath string, hard bool, revision string) error {
	cmd := NewCommand("reset")
	if hard {
		cmd.AddArguments("--hard")
	}
	_, err := cmd.AddArguments(revision).RunInDir(repoPath)
	return err
}