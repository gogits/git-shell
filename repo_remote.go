// Copyright 2019 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"strings"
	"time"
)

// LsRemoteOptions contains arguments for listing references in a remote repository.
// Docs: https://git-scm.com/docs/git-ls-remote
type LsRemoteOptions struct {
	// Indicates whether include heads.
	Heads bool
	// Indicates whether include tags.
	Tags bool
	// Indicates whether to not show peeled tags or pseudorefs.
	Refs bool
	// The list of patterns to filter results.
	Patterns []string
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// LsRemote returns a list references in the remote repository.
func LsRemote(url string, opts ...LsRemoteOptions) ([]*Reference, error) {
	var opt LsRemoteOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("ls-remote", "--quiet")
	if opt.Heads {
		cmd.AddArgs("--heads")
	}
	if opt.Tags {
		cmd.AddArgs("--tags")
	}
	if opt.Refs {
		cmd.AddArgs("--refs")
	}
	cmd.AddArgs(url)
	if len(opt.Patterns) > 0 {
		cmd.AddArgs(opt.Patterns...)
	}

	stdout, err := cmd.RunWithTimeout(opt.Timeout)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(stdout, []byte("\n"))
	refs := make([]*Reference, 0, len(lines))
	for i := range lines {
		fields := bytes.Fields(lines[i])
		if len(fields) < 2 {
			continue
		}

		refs = append(refs, &Reference{
			ID:      string(fields[0]),
			Refspec: string(fields[1]),
		})
	}
	return refs, nil
}

// IsURLAccessible returns true if given remote URL is accessible via Git
// within given timeout.
func IsURLAccessible(timeout time.Duration, url string) bool {
	_, err := LsRemote(url, LsRemoteOptions{
		Patterns: []string{"HEAD"},
		Timeout:  timeout,
	})
	return err == nil
}

// AddRemoteOptions contains options to add a remote address.
// Docs: https://git-scm.com/docs/git-remote#Documentation/git-remote.txt-emaddem
type AddRemoteOptions struct {
	// Indicates whether to execute git fetch after the remote information is set up.
	Fetch bool
	// Indicates whether to add remote as mirror with --mirror=fetch.
	MirrorFetch bool
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// AddRemote adds a new remote to the repository in given path.
func RepoAddRemote(repoPath, name, url string, opts ...AddRemoteOptions) error {
	var opt AddRemoteOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "add")
	if opt.Fetch {
		cmd.AddArgs("-f")
	}
	if opt.MirrorFetch {
		cmd.AddArgs("--mirror=fetch")
	}

	_, err := cmd.AddArgs(name, url).RunInDirWithTimeout(opt.Timeout, repoPath)
	return err
}

// AddRemote adds a new remote to the repository.
func (r *Repository) AddRemote(name, url string, opts ...AddRemoteOptions) error {
	return RepoAddRemote(r.path, name, url, opts...)
}

// RemoveRemoteOptions contains arguments for removing a remote from the repository.
// Docs: https://git-scm.com/docs/git-remote#Documentation/git-remote.txt-emremoveem
type RemoveRemoteOptions struct {
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// RepoRemoveRemote removes a remote from the repository in given path.
func RepoRemoveRemote(repoPath, name string, opts ...RemoveRemoteOptions) error {
	var opt RemoveRemoteOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	_, err := NewCommand("remote", "remove", name).RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil {
		// the error status may differ from git clients
		if strings.Contains(err.Error(), "error: No such remote") ||
			strings.Contains(err.Error(), "fatal: No such remote") {
			return ErrRemoteNotExist
		}
		return err
	}
	return nil
}

// RemoveRemote removes a remote from the repository.
func (r *Repository) RemoveRemote(name string, opts ...RemoveRemoteOptions) error {
	return RepoRemoveRemote(r.path, name, opts...)
}

// RemotesListOptions contains arguments for listing remotes of the repository.
// Docs: https://git-scm.com/docs/git-remote#_commands
type RemotesListOptions struct {
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// RepoRemotesList lists remotes of the repository in given path.
func RepoRemotesList(repoPath string, opts ...RemotesListOptions) ([]string, error) {
	var opt RemotesListOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	stdout, err := NewCommand("remote").RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil {
		return nil, err
	}

	return stdoutToStringSlice(stdout), nil
}

// RemotesList lists remotes of the repository.
func (r *Repository) RemotesList(opts ...RemotesListOptions) ([]string, error) {
	return RepoRemotesList(r.path, opts...)
}

// RemoteURLGetOptions contains arguments for retrieving URL(s) of a remote of the repository.
// Docs: https://git-scm.com/docs/git-remote#Documentation/git-remote.txt-emget-urlem
type RemoteURLGetOptions struct {
	// False: get fetch URLs
	// True: get push URLs
	Push bool
	// True: get all URLs (lists also non-main URLs; not related with Push)
	All bool
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// RepoRemoteURLGet retrieves URL(s) of a remote of the repository in given path.
func RepoRemoteURLGet(repoPath, name string, opts ...RemoteURLGetOptions) ([]string, error) {
	var opt RemoteURLGetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "get-url")
	if opt.Push {
		cmd.AddArgs("--push")
	}
	if opt.All {
		cmd.AddArgs("--all")
	}

	stdout, err := cmd.AddArgs(name).RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil {
		return nil, err
	}
	return stdoutToStringSlice(stdout), nil
}

// RemoteURLGet retrieves URL(s) of a remote of the repository in given path.
func (r *Repository) RemoteURLGet(name string, opts ...RemoteURLGetOptions) ([]string, error) {
	return RepoRemoteURLGet(r.path, name, opts...)
}

// RemoteURLSetOptions contains arguments for setting an URL of a remote of the repository.
// Docs: https://git-scm.com/docs/git-remote#Documentation/git-remote.txt-emget-urlem
type RemoteURLSetOptions struct {
	// False: set fetch URLs
	// True: set push URLs
	Push bool
	// The timeout duration before giving up for each shell command execution.
	// The default timeout duration will be used when not supplied.
	Timeout time.Duration
}

// RepoRemoteURLSetFirst sets first URL of the remote with given name of the repository in given path.
func RepoRemoteURLSetFirst(repoPath, name, newurl string, opts ...RemoteURLSetOptions) error {
	var opt RemoteURLSetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "set-url")
	if opt.Push {
		cmd.AddArgs("--push")
	}

	_, err := cmd.AddArgs(name, newurl).RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil && strings.Contains(err.Error(), "No such remote") {
		return ErrRemoteNotExist
	}
	return err
}

// RemoteURLSetFirst sets the first URL of the remote with given name of the repository.
func (r *Repository) RemoteURLSetFirst(name, newurl string, opts ...RemoteURLSetOptions) error {
	return RepoRemoteURLSetFirst(r.path, name, newurl, opts...)
}

// RepoRemoteURLSetRegex sets the first URL of the remote with given name and URL regex match of the repository in given path.
func RepoRemoteURLSetRegex(repoPath, name, urlregex string, newurl string, opts ...RemoteURLSetOptions) error {
	var opt RemoteURLSetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "set-url")
	if opt.Push {
		cmd.AddArgs("--push")
	}

	_, err := cmd.AddArgs(name, newurl, urlregex).RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil {
		if strings.Contains(err.Error(), "No such URL found") {
			return ErrURLNotExist
		}
		if strings.Contains(err.Error(), "No such remote") {
			return ErrRemoteNotExist
		}
		return err
	}
	return nil
}

// RemoteURLSetRegex sets the first URL of the remote with given name and URL regex match of the repository.
func (r *Repository) RemoteURLSetRegex(name, urlregex, newurl string, opts ...RemoteURLSetOptions) error {
	return RepoRemoteURLSetRegex(r.path, name, urlregex, newurl, opts...)
}

// RepoRemoteURLAdd adds an URL to the remote with given name of the repository in given path.
func RepoRemoteURLAdd(repoPath, name, newurl string, opts ...RemoteURLSetOptions) error {
	var opt RemoteURLSetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "set-url", "--add")
	if opt.Push {
		cmd.AddArgs("--push")
	}

	_, err := cmd.AddArgs(name, newurl).RunInDirWithTimeout(opt.Timeout, repoPath)
	return err
}

// RemoteURLAdd adds an URL to the remote with given name of the repository.
func (r *Repository) RemoteURLAdd(name, newvalue string, opts ...RemoteURLSetOptions) error {
	return RepoRemoteURLAdd(r.path, name, newvalue, opts...)
}

// RepoRemoteURLDelRegex Deletes all URLs matchin regex of the remote with given name of the repository in given path.
func RepoRemoteURLDelRegex(repoPath, name, urlregex string, opts ...RemoteURLSetOptions) error {
	var opt RemoteURLSetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	cmd := NewCommand("remote", "set-url", "--delete")
	if opt.Push {
		cmd.AddArgs("--push")
	}

	_, err := cmd.AddArgs(name, urlregex).RunInDirWithTimeout(opt.Timeout, repoPath)
	if err != nil && strings.Contains(err.Error(), "Will not delete all non-push URLs") {
		return ErrDelAllNonPushURL
	}
	return err
}

// RemoteURLDelRegex // RepoRemoteURLDelRegex Deletes all URLs matchin regex of the remote with given name of the repository.
func (r *Repository) RemoteURLDelRegex(name, urlregex string, opts ...RemoteURLSetOptions) error {
	return RepoRemoteURLDelRegex(r.path, name, urlregex, opts...)
}
