package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repository defines details to clone a dependency.
type Repository struct {
	// Dir contains the relative path where git will clone the repository from the root of the environment.
	Dir string      `json:"dir"`
	// URL contains the URL of the repository.
	URL string      `json:"url"`
	// Commit contains the hash that is checked out
	Commit string   `json:"commit"`
	// Commit contains the treeish (commit, tag or branch) that is used to select the commit
	Follow string   `json:"follow"`
}

// Build clones and checks out the repository at the specified commit inside the environment.
func (r *Repository) Build(root string) error {
	dir := filepath.Join(root, r.Dir)

	// git rev-list -n 1 --before="2017-12-03 00:00" master

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir

	msg := r.Dir

	if stdout, err := cmd.Output(); err == nil {
		if r.Commit == "" {
		  fmt.Println(msg, "-> at HEAD commit")
			return nil
		} else if r.Commit == strings.TrimSpace(string(stdout)) {
		  fmt.Println(msg, "-> at commit", r.Commit)
			return nil
		}

		// Already a cloned repo, just need to move to another commit
	} else {
		cmd = exec.Command("git", "clone", "-q", r.URL, r.Dir)
		cmd.Dir = root
		cmd.Stderr = os.Stderr
		cmd.Env = []string{
			"GIT_TERMINAL_PROMPT=0",
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s failed to 'git clone %s %s', got %s", msg, r.URL, r.Dir, err)
		}

		msg = fmt.Sprint(msg, " -> cloned")
	}

	if r.Commit == "" {
	  cmd = exec.Command("git", "rev-parse", "HEAD")
  	cmd.Dir = dir

  	if stdout, err := cmd.Output(); err == nil {
  		r.Commit = strings.TrimSpace(string(stdout))
		  fmt.Printf("%s -> checked out HEAD commit %s\n", msg, r.Commit)
  	} else {
			return fmt.Errorf("%s -> failed to 'git rev-parse HEAD', got %s", msg, r.URL, r.Dir, err)
		}

		return nil
	}

	cmd = exec.Command("git", "checkout", "-q", r.Commit)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s -> failed to 'git checkout %s', got %s", msg, r.Commit, err)
	}
	fmt.Printf("%s -> checked out commit %s\n", msg, r.Commit)

	return nil
}

//
func (r *Repository) Freeze(root string) error {
	dir := filepath.Join(root, r.Dir)

	// git rev-list -n 1 --before="2017-12-03 00:00" master

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir

	msg := r.Dir

	if stdout, err := cmd.Output(); err == nil {
	  commit := strings.TrimSpace(string(stdout))
		if r.Commit == "" {
		  msg = fmt.Sprint(msg, " -> frozen commit ",commit)
		} else if r.Commit == commit {
		  msg = fmt.Sprint(msg, " -> already frozen ")
		} else if r.Commit != commit {
		  msg = fmt.Sprint(msg, " -> upgraded frozen commit from ",r.Commit," to ",commit)
		}
		r.Commit = commit
	  fmt.Println(msg)
		// Already a cloned repo, just need to move to another commit
	} else {
		return fmt.Errorf("%s -> failed to 'git rev-parse HEAD', got %s", msg, err)
	}

	return nil
}
