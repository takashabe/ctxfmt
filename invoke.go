package main

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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

	dir, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: false,
		Dir:   dir,
	}
	pkg2, err := packages.Load(cfg, packageName)
	if err != nil {
		return err
	}
	pkg := pkg2[0].Errors
	if len(pkg) > 0 {
		return fmt.Errorf("packages.Load: %v", pkg)
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
