package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"unicode"
)

var ignoreFiles = []string{
	"mock_",
	"_sheet",
}

type inspectType string

const (
	inspectTypeInterface inspectType = "interface"
	inspectTypeStruct    inspectType = "method"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	fs := token.NewFileSet()
	for _, arg := range os.Args[2:] {
		file, err := parser.ParseFile(fs, arg, nil, 0)
		if err != nil {
			continue
		}

		if isIgnoreFile(fs.Position(file.Pos()).Filename) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name == nil {
				return true
			}

			if fn.Recv == nil {
				return false
			}

			if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
				return false
			}

			if !unicode.IsUpper(rune(fn.Name.Name[0])) {
				return false
			}

			if !hasContextParam(fn.Type.Params.List) {
				pos := fs.Position(fn.Pos())
				report(pos.Filename, fn.Name.Name, fn.Recv, pos.Line)
			}

			return false
		})
	}
}
