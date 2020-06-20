package lang

import (
	"testing"

	"github.com/iancoleman/strcase"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
		{"hello_svr", "hello_svr", "HelloSvr"},
		{"Hello_svr", "Hello_svr", "HelloSvr"},
		{"Hello_Svr", "Hello_Svr", "HelloSvr"},
		{"HelloSvr", "HelloSvr", "HelloSvr"},
		{"hellosvr", "hellosvr", "Hellosvr"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strcase.ToCamel(tt.args); got != tt.want {
				t.Errorf("strcase.ToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPBGoPackage(t *testing.T) {
	type args struct {
		pkgName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"case-1", args{"a.b.c"}, "a_b_c"},
		{"case-2", args{"a.b_c"}, "a_b_c"},
		{"case-3", args{"github.com/a/b.c_d"}, "github.com/a/b_c_d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PBGoPackage(tt.args.pkgName); got != tt.want {
				t.Errorf("PBGoPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
