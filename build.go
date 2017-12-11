package main

import (
	"flag"
	"os"
)

func build(path string, args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	jobs := flags.Int("jobs", 4, "number of concurrent git clone jobs")
	upgrade := flags.Bool("upgrade", false, "upgrades dependencies")
	reset := flags.Bool("reset", false, "clears the .env directory prior to build")
	timestamp := flags.String("timestamp", "", "use commits at that point in time")
	flags.Parse(args)

	env, err := New(path)
	if err != nil {
		return err
	}

	if *upgrade && *timestamp != "" {
		// This is an error - can't upgrade AND sync to a given timestamp
	}

	if reset != nil && *reset {
		if err := os.RemoveAll(env.path); err != nil {
			return err
		}
	}

	env.Jobs = *jobs

	return env.Build(*upgrade, *timestamp)
}
