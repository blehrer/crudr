package util

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Ast() {
	fset := token.NewFileSet()
	goPath := os.Getenv("GOPATH")
	pkgPath := goPath + "/pkg/mod/github.com/pb33f/libopenapi*"

	matches, _ := filepath.Glob(pkgPath)
	match := matches[0]
	var dirs []string
	err := filepath.Walk(match, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() &&
			!strings.Contains(path, "test") &&
			!strings.Contains(path, "low") {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, func(info fs.FileInfo) bool {
			return !strings.Contains(info.Name(), "test")
		}, 0)
		if err != nil {
			log.Fatal(err)
		}
		//for name, pkg := range pkgs {
		for _, pkg := range pkgs {
			if true {
				//if slices.Contains([]string{"v3", "base", "model"}, name) {
				//if slices.Contains([]string{"base"}, name) {
				//if slices.Contains([]string{"model"}, name) {
				//if slices.Contains([]string{"v3"}, name) {
				fmt.Printf("%s\n", pkg.Name)
				for _, file := range pkg.Files {
					// Process each file's AST
					// Inside the loop processing each file's AST
					ast.Inspect(file, func(n ast.Node) bool {
						if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
							for _, spec := range genDecl.Specs {
								if tspec, ok := spec.(*ast.TypeSpec); ok {
									if _, isStruct := tspec.Type.(*ast.StructType); isStruct {
										// Found a struct declaration
										structName := tspec.Name.Name
										// filter out private structs
										firstChar := string(structName[0])
										if firstChar == strings.ToUpper(firstChar) {
											fmt.Printf("- %s\n", structName)
										}
									}
								}
							}
						}
						return true // Continue traversing the AST
					})
				}
				fmt.Println("")
			}
		}
	}
}
