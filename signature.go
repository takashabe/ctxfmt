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

func fmtSignature(fs *token.FileSet, fileName string, dryrun bool) error {
	info, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(fileName, ".go") {
		return nil
	}

	file, err := parser.ParseFile(fs, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astutil.AddImport(fs, file, "context")
	var isApply bool
	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		switch decl := cr.Node().(type) {
		case *ast.FuncDecl:
			if decl.Name != nil {
				if !isIgnoreFunc(decl.Name.Name) {
					if decl.Recv != nil {
						if decl.Type.Params != nil && len(decl.Type.Params.List) > 0 {
							if unicode.IsUpper(rune(decl.Name.Name[0])) {
								if !hasContextParam(decl.Type.Params.List) {
									if dryrun {
										pos := fs.Position(decl.Pos())
										reportMethod(pos.Filename, decl.Name.Name, decl.Recv, pos.Line)
									} else {
										contextParam := &ast.Field{
											Names: []*ast.Ident{ast.NewIdent("ctx")},
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("context"),
												Sel: ast.NewIdent("Context"),
											},
										}
										decl.Type.Params.List = append([]*ast.Field{contextParam}, decl.Type.Params.List...)
										isApply = true
									}
								}
							}
						}
					}
				}
			}
		case *ast.TypeSpec:
			if interfaceType, ok := decl.Type.(*ast.InterfaceType); ok {
				for _, m := range interfaceType.Methods.List {
					if method, ok := m.Type.(*ast.FuncType); ok {
						if !hasContextParam(method.Params.List) {
							if dryrun {
								pos := fs.Position(m.Pos())
								reportInterface(pos.Filename, m.Names[0].Name, decl, pos.Line)
							} else {
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
		}

		return true
	}, nil)

	if !isApply {
		return nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fs, file); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("processed %s\n", fileName)

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
