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
	configFileFlag = &cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "config file path",
		Destination: &configFile,
	}
)

var (
	skipDefinedMethod     bool
	skipDefinedMethodFlag = &cli.BoolFlag{
		Name:        "skip-method",
		Usage:       "skip defined method (not interface)",
		Destination: &skipDefinedMethod,
	}
)

var (
	skipInterface     bool
	skipInterfaceFlag = &cli.BoolFlag{
		Name:        "skip-interface",
		Usage:       "skip declared interface method",
		Destination: &skipInterface,
	}
)

func main() {
	app := &cli.App{
		Name:     "ctxfmt",
		HelpName: "",
		Usage:    "context.Context formatter",
		Commands: []*cli.Command{
			{
				Name:  "def",
				Usage: "format method definition",
				Flags: []cli.Flag{
					dryrunFlag,
					configFileFlag,
					skipDefinedMethodFlag,
					skipInterfaceFlag,
				},
				ArgsUsage: "target file or directory",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}

					if err := loadConfig(configFile); err != nil {
						return fmt.Errorf("failed to load config: %w", err)
					}

					fs := token.NewFileSet()
					for _, arg := range c.Args().Slice() {
						if err := fmtDef(fs, arg, formatDefConfig{
							IgnoreFuncs:   ignoreFuncs,
							Dryrun:        dryrun,
							SkipMethod:    skipDefinedMethod,
							SkipInterface: skipInterface,
						}); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:  "call",
				Usage: "format method call",
				Flags: []cli.Flag{
					dryrunFlag,
					pkgFlag,
					configFileFlag,
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
						if err := fmtCall(fs, arg, pkgName, dryrun); err != nil {
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

func isIgnoreFunc(target string, ignores []string) bool {
	for _, ig := range ignores {
		if target == ig {
			return true
		}
	}
	return false
}
