package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var ignoreFuncs = []string{
	"PreInsert",
	"PreUpdate",
	"Scan",
}

var dryrun bool

var dryrunFlag = &cli.BoolFlag{
	Name:        "dryrun",
	Aliases:     []string{"n"},
	Usage:       "dryrun",
	Destination: &dryrun,
}

var pkgName string

var pkgFlag = &cli.StringFlag{
	Name:        "pkg",
	Aliases:     []string{"p"},
	Usage:       "package name",
	Destination: &pkgName,
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "signature",
				Usage: "format function / method signature",
				Flags: []cli.Flag{
					dryrunFlag,
					pkgFlag,
				},
				ArgsUsage: "file/dir",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}
					if pkgName == "" {
						return fmt.Errorf("invalid pkg name")
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
				ArgsUsage: "file/dir",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}
					if pkgName == "" {
						return fmt.Errorf("invalid pkg name")
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
