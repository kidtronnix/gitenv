package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// New creates an empty environment or opens it if it already exists.
func New(path string) (env *Environment, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}

	dir := filepath.Dir(path)

	env = &Environment{
		file: path,
		root: dir,
		path: filepath.Join(dir, ".env"),
	}

	err = env.Load()
	return
}

// Environment defines the details required to manage a repository.
type Environment struct {
	// Links maps source directories and files in the environment.
	Links map[string]string
	// Repositories contains repositories to clone.
	Repositories []*Repository

	file string
	root string
	path string
}

// Load fills the environment's details from the .gitenv file.
func (env *Environment) Load() error {
	file := env.file

	if file == "" {
		if env.root != "" {
			file = filepath.Join(env.root, ".gitenv")
		} else {
			file = "./.gitenv"
		}
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read: %s", err)
	}

	if err := json.Unmarshal(content, env); err != nil {
		return fmt.Errorf("json: %s", err)
	}

	return nil
}

// Build sets links and clone repositories in the environment's directory.
func (env *Environment) Build() error {
	fmt.Println("build:", env.path)

	symlink := func(source, target string) error {
		if a, err := os.Stat(target); err == nil {
			if b, err := os.Stat(source); err == nil {
				if os.SameFile(a, b) {
					return nil
				}
			} else {
				return err
			}

			return fmt.Errorf("%s: already exists", target)
		}

		if os.Symlink(source, target) == nil {
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(target), 0777); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}

		return os.Symlink(source, target)
	}

	for link, path := range env.Links {
		source := filepath.Join(env.root, link)
		target := filepath.Join(env.path, path)
		if err := symlink(source, target); err != nil {
			return err
		}
	}

	n := len(env.Repositories)
	if n == 0 {
		return nil
	}

	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		r := env.Repositories[i]
		go func() {
			errs <- r.Build(env.path)
		}()
	}

	fails := 0

	for i := 0; i < n; i++ {
		if err := <-errs; err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			fails++
		}
	}

	if fails != 0 {
		return fmt.Errorf("%d errors", fails)
	}

	return nil
}
