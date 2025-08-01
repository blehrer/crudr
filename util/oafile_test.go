package util

import (
	"github.com/charmbracelet/huh"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"reflect"
	"testing"
)

func TestToForm(t *testing.T) {
	type args struct {
		datastructure any
	}
	tests := []struct {
		name string
		args args
		want map[string]func(s string) huh.Group
	}{
		// TODO: Add test cases.
		{
			name: "ToForm-libopenapi.v3.PathItem",
			args: args{
				datastructure: v3.PathItem{
					Description: "",
					Summary:     "",
					Get:         nil,
					Put:         nil,
					Post:        nil,
					Delete:      nil,
					Options:     nil,
					Head:        nil,
					Patch:       nil,
					Trace:       nil,
					Servers:     nil,
					Parameters:  nil,
					Extensions:  nil,
				},
			},
			want: func() map[string]func(s string) huh.Group {
				t := reflect.TypeFor[v3.PathItem]()
				rv := make(map[string]func(s string) huh.Group, t.NumField())
				rv["Description"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Description"))
				}
				rv["Summary"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Summary"))
				}
				rv["Get"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Get"))
				}
				rv["Put"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Put"))
				}
				rv["Post"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Post"))
				}
				rv["Delete"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Delete"))
				}
				rv["Options"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Options"))
				}
				rv["Head"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Head"))
				}
				rv["Patch"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Patch"))
				}
				rv["Trace"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Trace"))
				}
				rv["Servers"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Servers"))
				}
				rv["Parameters"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Parameters"))
				}
				rv["Extensions"] = func(s string) huh.Group {
					return *huh.NewGroup(
						huh.NewInput().
							Title("Extensions"))
				}
				return rv
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToForm(tt.args.datastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToForm() = %v, want %v", got, tt.want)
			}
		})
	}
}
