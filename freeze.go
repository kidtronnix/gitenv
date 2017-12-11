package main

import (
	"flag"
)

func freeze(path string, args []string) error {
	flags := flag.NewFlagSet("freeze", flag.ExitOnError)
	// reset := flags.Bool("reset", false, "clears the .env directory prior to build")
	flags.Parse(args)

	env, err := New(path)
	if err != nil {
		return err
	}

	if err = env.Freeze(); err != nil {
		return err
	}

	return env.Save()
}
