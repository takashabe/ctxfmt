package main

import (
	"fmt"
	"go/token"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var ignoreFiles = []string{
	"mock_",
	"_sheet",
}

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

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:      "interface",
				Aliases:   []string{"i"},
				Usage:     "format interface",
				Flags:     []cli.Flag{dryrunFlag},
				ArgsUsage: "file/dir",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}

					fs := token.NewFileSet()
					for _, arg := range c.Args().Slice() {
						if err := fmtInterface(fs, arg, dryrun); err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:      "method",
				Aliases:   []string{"m"},
				Usage:     "format method",
				Flags:     []cli.Flag{dryrunFlag},
				ArgsUsage: "file/dir",
				Action: func(c *cli.Context) error {
					if c.Args().Len() < 1 {
						return fmt.Errorf("invalid args")
					}

					fs := token.NewFileSet()
					for _, arg := range c.Args().Slice() {
						if err := fmtMethod(fs, arg, dryrun); err != nil {
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
