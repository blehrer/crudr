package util

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

var drilldowns map[string]string

// Dump parses the abstract syntax tree of a given lib, for a given set of its
// packages. Currently, it prints to stdout, but eventually this should be used
// to seed a command that allows you to pick forms for each type of struct.
// libUrl should be in the same format as what you would import (e.g.
// "github.com/author/lib"). You can leave version empty to default to "*".
func Dump(libUrl string, version string, packages []string) {
	drilldowns = map[string]string{}
	fset := token.NewFileSet()
	goPath := os.Getenv("GOPATH")
	v := "*"
	if version != "" {
		v = version
	}
	pkgPath := fmt.Sprintf("%s/pkg/mod/%s%s", goPath, libUrl, v)

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
		for name, pkg := range pkgs {
			if slices.Contains(packages, name) {
				for _, file := range pkg.Files {
					// Process each file's AST
					// Inside the loop processing each file's AST
					ast.Inspect(file, func(n ast.Node) bool {
						drilldown([]ast.Node{n})
						return true // Continue traversing the AST
					})
				}
			}
		}
	}
	sortedKeys := []string{}
	for k := range maps.Keys(drilldowns) {
		sortedKeys = append(sortedKeys, k)
	}
	for k := range sortedKeys {
		fmt.Println(sortedKeys[k])
		fmt.Println(drilldowns[sortedKeys[k]])
	}
}

func drilldown(ns []ast.Node) []ast.Node {
	if len(ns) == 0 {
		return nil
	}
	var nextNodes []ast.Node
	for _, n := range ns {
		// filter inspection to token.TYPE
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			specs := n.(*ast.GenDecl).Specs
			if len(specs) != 1 {
				log.Fatal()
			}
			for _, spec := range genDecl.Specs {
				if tspec, ok := spec.(*ast.TypeSpec); ok {
					if _, isStruct := tspec.Type.(*ast.StructType); isStruct && tspec.Name.IsExported() {
						// Found a struct declaration
						structName := tspec.Name.String()
						fields := tspec.Type.(*ast.StructType).Fields.List
						var fieldDescriptions string
						for _, field := range fields {
							fieldName := field.Names[0].Name
							if string(fieldName[0]) == strings.ToUpper(string(fieldName[0])) {
								typeNameIface := reflect.ValueOf(field.Type).Interface()
								var typeName string
								typeName = formatTypeInfo(typeNameIface, typeName)
								fieldDescriptions = fieldDescriptions + fmt.Sprintf("\t- %-20s\t%s\n",
									fieldName, typeName)
								nextNodes = append(nextNodes, ast.NewIdent(fieldName))
							}
						}

						// filter out private structs
						firstChar := string(structName[0])
						if firstChar == strings.ToUpper(firstChar) {
							drilldowns[structName] = fieldDescriptions
						}
					}
				}
			}
		}

	}
	return drilldown(nextNodes)
}

func formatTypeInfo(typeNameIface any, typeName string) string {
	switch ftype := typeNameIface.(type) {
	case *ast.Ident:
		typeName = ftype.Name
	case *ast.ArrayType:
		typeName = fmt.Sprintf("[]%s", formatTypeInfo(ftype.Elt, typeName))
	case *ast.MapType:
		var keytype string
		var valtype string
		formatTypeInfo(ftype.Key, keytype)
		formatTypeInfo(ftype.Value, valtype)
		typeName = fmt.Sprintf("map[%s]%s", keytype, valtype)
	case *ast.IndexListExpr:
		var indexTypes []string
		for _, innerFtype := range ftype.Indices {
			var innerTypeName string
			indexTypes = append(indexTypes, formatTypeInfo(innerFtype, innerTypeName))
		}
		var indexedType string
		indexedType = formatTypeInfo(ftype.X, indexedType)
		typeName = fmt.Sprintf("%s[%v]", indexedType, strings.Join(indexTypes, ", "))
	case *ast.StarExpr:
		typeName = fmt.Sprintf("*%s", formatTypeInfo(ftype.X, typeName))
	case *ast.SelectorExpr:
		if ftype.Sel.IsExported() {
			// todo: look up ftype.X in ast, and find the lib/pkg it refers to.
			typeName = fmt.Sprintf("%v.%v", formatTypeInfo(ftype.X, typeName), ftype.Sel.String())
		}
	default:
		log.Fatal(ftype)
	}
	return typeName
}
