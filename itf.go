package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"unicode"
)

func inspectInterface(fs *token.FileSet, filename string) error {
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		for _, m := range interfaceType.Methods.List {
			for _, name := range m.Names {
				if !hasContextParam(m.Type.(*ast.FuncType).Params.List) {
					fmt.Printf("In file %s: Method %s of interface %s does not take context.Context\n", filename, name.Name, typeSpec.Name.Name)
				}
			}
		}
		return true
	})

	return nil
}

// interface用のreport作る
func reportInterface(filename, funcName string, typeSpec *ast.TypeSpec, line int) {
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

func inspectMethod(fs *token.FileSet, filename string) error {
	file, err := parser.ParseFile(fs, filename, nil, 0)
	if err != nil {
		return err
	}

	if isIgnoreFile(fs.Position(file.Pos()).Filename) {
		return nil
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
			reportMethod(pos.Filename, fn.Name.Name, fn.Recv, pos.Line)
		}

		return false
	})
	return nil
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

func reportMethod(filename, funcName string, recv *ast.FieldList, line int) {
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
