package main

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func invoke(fs *token.FileSet, filename string, packageName string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		// skip traversal
		return nil
	}

	pkgs, err := parser.ParseDir(fs, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return nil
	}

	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: false,
	}
	pkgs2, err := packages.Load(cfg, packageName)
	if err != nil {
		return err
	}

	// パッケージの解析
	pkg := pkgs2[0]
	if len(pkg.Errors) > 0 {
		return fmt.Errorf("packages.Load: %v", pkg.Errors)
	}

	return nil
}

// extractModuleName はgo.modファイルからモジュール名を抽出します。
func extractModuleName(goModPath string) (string, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("module directive not found in %s", goModPath)
}

// package main
//
// import (
//   "go/ast"
//   "go/parser"
//   "go/token"
//   "go/types"
//   "log"
//   "os"
// )
//
// func main() {
//   filename := "yourfile.go" // 解析するファイル
//
//   // ファイルセットの作成
//   fset := token.NewFileSet()
//
//   // ソースファイルの解析
//   file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
//   if err != nil {
//     log.Fatal(err)
//   }
//
//   // 型チェッカーの設定
//   conf := types.Config{Importer: types.DefaultImporter}
//   info := &types.Info{
//     Types: make(map[ast.Expr]types.TypeAndValue),
//   }
//
//   // ファイルに含まれるパッケージの型情報を取得
//   pkg, err := conf.Check("main", fset, []*ast.File{file}, info)
//   if err != nil {
//     log.Fatal(err)
//   }
//
//   // ASTをトラバースし、関数呼び出しを探す
//   ast.Inspect(file, func(n ast.Node) bool {
//     if callExpr, ok := n.(*ast.CallExpr); ok {
//       fnType, ok := info.Types[callExpr.Fun].Type.(*types.Signature)
//       if !ok {
//         return true
//       }
//
//       // 関数の引数が足りているかチェック
//       if len(callExpr.Args) < fnType.Params().Len() {
//         // context.Context が不足しているか確認
//         // ...
//       }
//     }
//     return true
//   })
// }
