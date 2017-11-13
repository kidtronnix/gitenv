package main

import (
	"flag"
	"os"
)

func build(path string, args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	jobs := flags.Int("jobs", 4, "number of concurrent git clone jobs")
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

	env.Jobs = *jobs

	return env.Build()
}
