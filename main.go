package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v2"
)

// ignoreFuncs is a list of function names to ignore.
var ignoreFuncs []string

var (
	dryrun     bool
	dryrunFlag = &cli.BoolFlag{
		Name:        "dryrun",
		Aliases:     []string{"n"},
		Usage:       "dryrun",
		Destination: &dryrun,
	}
)

var (
	pkgName string
	pkgFlag = &cli.StringFlag{
		Name:        "pkg",
		Aliases:     []string{"p"},
		Usage:       "package name",
		Destination: &pkgName,
	}
)

type config struct {
	IgnoreFuncs []string `yaml:"ignore_funcs"`
}

var (
	configFile     string
	configFilePath = &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "config file path",
		Destination: &configFile,
	}
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "signature",
				Usage: "format function, method signature",
				Flags: []cli.Flag{
					dryrunFlag,
					pkgFlag,
				},
				ArgsUsage: "target file or directory",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}
					if pkgName == "" {
						return fmt.Errorf("invalid pkg name")
					}

					if err := loadConfig(configFile); err != nil {
						return fmt.Errorf("failed to load config: %w", err)
					}

					fs := token.NewFileSet()
					for _, arg := range c.Args().Slice() {
						if err := fmtSignature(fs, arg, dryrun); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name: "arg",
				Flags: []cli.Flag{
					dryrunFlag,
					pkgFlag,
				},
				ArgsUsage: "target directory",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}
					if pkgName == "" {
						return fmt.Errorf("invalid pkg name")
					}

					if err := loadConfig(configFile); err != nil {
						return fmt.Errorf("failed to load config: %w", err)
					}

					fs := token.NewFileSet()
					for _, arg := range c.Args().Slice() {
						if err := fmtArgs(fs, arg, pkgName, dryrun); err != nil {
							return err
						}
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func loadConfig(configFile string) error {
	if configFile == "" {
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var cfg config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return err
	}

	// TODO: support command line args
	if len(cfg.IgnoreFuncs) > 0 {
		ignoreFuncs = cfg.IgnoreFuncs
	}

	return nil
}

func isIgnoreFunc(name string) bool {
	for _, ignoreFunc := range ignoreFuncs {
		if name == ignoreFunc {
			return true
		}
	}
	return false
}
