package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func fmtArgs(fs *token.FileSet, filename, pkgName string, dryrun bool) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return nil
	}

	pkgDir, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  pkgDir,
	}

	pkgs, err := packages.Load(cfg, pkgName)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			funcNames, ok := notEnoughContextArgs(err.Error())
			if !ok {
				continue
			}
			for _, name := range funcNames {
				for _, file := range pkg.CompiledGoFiles {
					if err := addContextToFunctionCall(fs, file, name); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func addContextToFunctionCall(fs *token.FileSet, fileName, funcName string) error {
	file, err := parser.ParseFile(fs, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	astutil.AddImport(fs, file, "context")
	var isApply bool
	astutil.Apply(file, func(cr *astutil.Cursor) bool {
		if callExpr, ok := cr.Node().(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok && selectorExpr.Sel.Name == funcName {
					// 既に最初の引数が context.Context のようなものかどうかをチェック
					if len(callExpr.Args) > 0 {
						firstArg := callExpr.Args[0]
						switch arg := firstArg.(type) {
						case *ast.Ident:
							if arg.Name == "ctx" {
								// 最初の引数が "ctx" なのでスキップ
								return true
							}
						case *ast.SelectorExpr:
							// 第一引数がcontextパッケージの関数呼び出しであればスキップ
							if xIdent, ok := arg.X.(*ast.Ident); ok && xIdent.Name == "context" {
								return true
							}
						}
					}

					if dryrun {
						pos := fs.Position(ident.Pos())
						reportArgs(pos.Filename, funcName, pos.Line)
					} else {
						contextCall := ast.NewIdent("context.TODO()")
						callExpr.Args = append([]ast.Expr{contextCall}, callExpr.Args...)
						isApply = true
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

func reportArgs(filename, funcName string, line int) {
	fmt.Printf("%s at line %d: %s()\n", filename, line, funcName)
}

var (
	interfaceMethodRegex = regexp.MustCompile(`want (\w+)\(context\.Context,`)
	functionCallRegex    = regexp.MustCompile(`not enough arguments in call to [\w.]+\b\.(\w+)`)
)

// notEnoughContextArgs returns function names that have not enough context.Context arguments.
func notEnoughContextArgs(errMessage string) ([]string, bool) {
	var funcNames []string

	for _, ptn := range []*regexp.Regexp{interfaceMethodRegex, functionCallRegex} {
		matches := ptn.FindAllStringSubmatch(errMessage, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				funcNames = append(funcNames, match[1])
			}
		}
	}
	return funcNames, len(funcNames) > 0
}
