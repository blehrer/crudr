package util

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func Test_main(t *testing.T) {
	type TypeA struct {
		aName string
	}
	type TypeB struct {
		bName string
		aType TypeA
	}
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Dump("github.com/pb33f/libopenapi", "", []string{"v3"})
			rb := reflect.TypeFor[TypeB]()
			rbFields := make([]reflect.StructField, rb.NumField())
			for i := 0; i < rb.NumField(); i++ {
				rbFields[i] = rb.Field(i)
				var pos token.Pos
				var pkg *types.Package
				var name string
				var typ types.Type
				types.NewVar(pos, pkg, name, typ)
			}
		})
	}
}
