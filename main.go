package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
)

var ignoreFiles = []string{
	"mock_",
	"_sheet",
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	fs := token.NewFileSet()
	for _, arg := range os.Args[1:] {
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

func hasContextParam(fields []*ast.Field) bool {
	for _, field := range fields {
		if typExpr, ok := field.Type.(*ast.SelectorExpr); ok {
			if pkgIdent, ok := typExpr.X.(*ast.Ident); ok {
				if pkgIdent.Name == "context" && typExpr.Sel.Name == "Context" {
					return true
				}
			}
		}
	}
	return false
}

func report(filename, funcName string, recv *ast.FieldList, line int) {
	var recvName string
	if recv != nil && len(recv.List) > 0 {
		if len(recv.List[0].Names) > 0 {
			recvName = recv.List[0].Names[0].Name
		} else if recvType, ok := recv.List[0].Type.(*ast.Ident); ok {
			recvName = recvType.Name
		}
	}

	fmt.Printf("%s at line %d: %s.%s()\n", filename, line, recvName, funcName)
}

func isIgnoreFile(fileName string) bool {
	for _, ignoreFile := range ignoreFiles {
		if strings.Contains(fileName, ignoreFile) {
			return true
		}
	}
	return false
}
