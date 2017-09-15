package main

import (
	"flag"
	"os"
)

func build(path string, args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	reset := flags.Bool("reset", false, "clears the .env directory prior to build")
	flags.Parse(args)

	env, err := New(path)
	if err != nil {
		return err
	}

	if reset != nil && *reset {
		if err := os.RemoveAll(env.path); err != nil {
			return err
		}
	}

	return env.Build()
}
