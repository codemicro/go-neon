package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/codemicro/go-neon/neontc/config"
	"github.com/codemicro/go-neon/neontc/tc"
	"os"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unhandled error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	conf := new(config.Config)
	arg.MustParse(conf)

	if len(conf.Directory) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		conf.Directory = cwd
	}

	if len(conf.OutputDirectory) == 0 {
		conf.OutputDirectory = conf.Directory
	}

	if err := tc.RunOnDirectory(conf, conf.Directory); err != nil {
		return err
	}

	return nil
}
