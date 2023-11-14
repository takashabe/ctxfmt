package main

import (
	"go/token"
	"os"
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

type inspectType string

const (
	inspectTypeInterface  inspectType = "interface"
	inspectTypeInterface2 inspectType = "rewrite_interface"
	inspectTypeMethod     inspectType = "method"
	inspectTypeMethod2    inspectType = "rewrite_method"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	switch typ := os.Args[1]; typ {
	case string(inspectTypeInterface):
		fs := token.NewFileSet()
		for _, arg := range os.Args[2:] {
			if err := inspectInterface(fs, arg); err != nil {
				panic(err)
			}
		}
	case string(inspectTypeInterface2):
		fs := token.NewFileSet()
		for _, arg := range os.Args[2:] {
			if err := rewirteInterface(fs, arg); err != nil {
				panic(err)
			}
		}
	case string(inspectTypeMethod):
		fs := token.NewFileSet()
		for _, arg := range os.Args[2:] {
			if err := inspectMethod(fs, arg); err != nil {
				panic(err)
			}
		}
	case string(inspectTypeMethod2):
		fs := token.NewFileSet()
		for _, arg := range os.Args[2:] {
			if err := inspectMethod2(fs, arg); err != nil {
				panic(err)
			}
		}
	default:
		panic("invalid type")
	}
}
