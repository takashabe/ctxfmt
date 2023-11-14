package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"

	"golang.org/x/tools/go/ast/astutil"
)

func inspectInterface(fs *token.FileSet, filename string) error {
	file, err := parser.ParseFile(fs, filename, nil, 0)
	if err != nil {
		return nil
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				for _, m := range interfaceType.Methods.List {
					for _, name := range m.Names {
						if !hasContextParam(m.Type.(*ast.FuncType).Params.List) {
							pos := fs.Position(m.Pos())
							reportInterface(pos.Filename, name.Name, typeSpec, pos.Line)
						}
					}
				}
			}
		}
		return true
	})

	return nil
}

func rewirteInterface(fs *token.FileSet, filename string) error {
	file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astutil.AddImport(fs, file, "context")
	var isApply bool
	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		n := cr.Node()
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				for _, m := range interfaceType.Methods.List {
					if method, ok := m.Type.(*ast.FuncType); ok {
						if !hasContextParam(method.Params.List) {
							contextParam := &ast.Field{
								Names: []*ast.Ident{ast.NewIdent("ctx")},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
								},
							}
							method.Params.List = append([]*ast.Field{contextParam}, method.Params.List...)
							isApply = true
						}
					}
				}
			}
		}
		return true
	}, nil)

	if !isApply {
		return nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fs, file); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("processed %s\n", filename)

	return nil
}

func inspectMethod(fs *token.FileSet, filename string) error {
	file, err := parser.ParseFile(fs, filename, nil, 0)
	if err != nil {
		return nil
	}

	if isIgnoreFile(fs.Position(file.Pos()).Filename) {
		return nil
	}

	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name != nil {
			if !isIgnoreFuncName(fn.Name.Name) {
				if fn.Recv != nil {
					if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
						if unicode.IsUpper(rune(fn.Name.Name[0])) {
							if !hasContextParam(fn.Type.Params.List) {
								pos := fs.Position(fn.Pos())
								reportMethod(pos.Filename, fn.Name.Name, fn.Recv, pos.Line)
							}
						}
					}
				}
			}
		}

		return true
	})
	return nil
}

func inspectMethod2(fs *token.FileSet, filename string) error {
	file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		return nil
	}
	if isIgnoreFile(fs.Position(file.Pos()).Filename) {
		return nil
	}

	astutil.AddImport(fs, file, "context")
	var isApply bool
	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		n := cr.Node()
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name != nil {
			if !isIgnoreFuncName(fn.Name.Name) {
				if fn.Recv != nil {
					if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
						if unicode.IsUpper(rune(fn.Name.Name[0])) {
							if !hasContextParam(fn.Type.Params.List) {
								contextParam := &ast.Field{
									Names: []*ast.Ident{ast.NewIdent("ctx")},
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("context"),
										Sel: ast.NewIdent("Context"),
									},
								}
								fn.Type.Params.List = append([]*ast.Field{contextParam}, fn.Type.Params.List...)
								isApply = true
							}
						}
					}
				}
			}
		}
		return true
	}, nil)

	if !isApply {
		return nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fs, file); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("processed %s\n", filename)

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

func reportInterface(filename, funcName string, typeSpec *ast.TypeSpec, line int) {
	fmt.Printf("%s at line %d: %s.%s()\n", filename, line, typeSpec.Name.Name, funcName)
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

func isIgnoreFuncName(funcName string) bool {
	for _, ignoreFunc := range ignoreFuncs {
		if strings.Contains(funcName, ignoreFunc) {
			return true
		}
	}
	return false
}
