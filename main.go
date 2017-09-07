package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gitenv: %s\n", err)
		os.Exit(1)
	}

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	path := flags.String("path", filepath.Join(pwd, ".gitenv"), "Location of the .gitenv `file`.")

	flags.Usage = func() {
		fmt.Println("Usage of gitenv: [flags...] command [args...]")
		fmt.Println("")
		fmt.Println("commands:")
		fmt.Println("")
		fmt.Println("  build       create/update .env from .gitenv file")
		fmt.Println("  help        display help on a command")
		fmt.Println("")
		fmt.Println("flags:")
		fmt.Println("")
		flags.PrintDefaults()
		os.Exit(2)
	}

	flags.Parse(os.Args[1:])

	cmd, args := flags.Arg(0), flags.Args()

	if len(args) == 0 {
		flags.Usage()
	}

	if cmd == "help" {
		if len(args) == 1 {
			flags.Usage()
		}

		cmd, args = args[1], []string{"--help"}
	} else {
		args = args[1:]
	}

	run := func(err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "gitenv %s: %s\n", cmd, err)
			os.Exit(1)
		}
	}

	switch cmd {
	case "build":
		run(build(*path, args))
	default:
		fmt.Fprintf(os.Stderr, "gitenv: '%s' is not a command. See 'gitenv --help'\n", cmd)
		os.Exit(2)
	}
}
