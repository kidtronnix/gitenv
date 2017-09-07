package main

import (
	"flag"
)

func build(path string, args []string) error {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	flags.Parse(args)

	env, err := New(path)
	if err != nil {
		return err
	}

	return env.Build()
}
