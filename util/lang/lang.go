package lang

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
)

// PBSimplifyGoType determine whether to use fullyQualifiedPackageName or not,
// if the `fullTypeName` occur in code of `package goPackageName`, `package` part
// should be removed.
func PBSimplifyGoType(fullTypeName string, goPackageName string) string {

	idx := strings.LastIndex(fullTypeName, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("invalid fullyQualifiedType:%s", fullTypeName))
	}

	pkg := fullTypeName[0:idx]
	typ := fullTypeName[idx+1:]

	if pkg == goPackageName {
		//fmt.Println("pkg:", pkg, "=", "gopkg:", goPackageName)
		return typ
	}
	//fmt.Println("pkg:", pkg, "!=", "gopkg:", goPackageName)
	return fullTypeName
}

// PBGoType convert `t` to go style (like a.b.c.hello, it'll be changed to a_b_c.Hello)
func PBGoType(t string) string {

	var prefix string

	idx := strings.LastIndex(t, "/")
	if idx >= 0 {
		prefix = t[:idx]
		t = t[idx+1:]
	}

	idx = strings.LastIndex(t, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("fatal error: invalid type:%s", t))
	}

	gopkg := PBGoPackage(t[0:idx])
	msg := t[idx+1:]

	return GoExport(prefix + gopkg + "." + msg)
}

// PBGoPackage convert a.b.c to a_b_c
func PBGoPackage(pkgName string) string {
	var (
		prefix string
		pkg    string
	)
	idx := strings.LastIndex(pkgName, "/")
	if idx < 0 {
		pkg = pkgName
	} else {
		prefix = pkgName[0:idx]
		pkg = pkgName[idx+1:]
	}

	gopkg := strings.Replace(pkg, ".", "_", -1)

	if len(prefix) == 0 {
		return gopkg
	}
	return prefix + "/" + gopkg
}

// GoExport export go type
func GoExport(typ string) string {
	idx := strings.LastIndex(typ, ".")
	if idx < 0 {
		return strings.Title(typ)
	}
	return typ[0:idx] + "." + strings.Title(typ[idx+1:])
}

// SplitList split string `str` via delimiter `sep` into a list of string
func SplitList(sep, str string) []string {
	return strings.Split(str, sep)
}

// TrimRight trim right substr starting at `sep`
func TrimRight(sep, str string) string {
	idx := strings.LastIndex(str, sep)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

// TrimLeft trim left substr starting at `sep`
func TrimLeft(sep, str string) string {
	return strings.TrimPrefix(str, sep)
}

// Title uppercase the first character of `s`
func Title(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

func UnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

// GoFullyQualifiedType convert $repo/$pkg.$type to $realpkg.$type, where $realpkg is calculated
// by `package directive` and `go_package` file option.
func GoFullyQualifiedType(pbFullyQualifiedType string, nfd *descriptor.FileDescriptor) string {

	idx := strings.LastIndex(pbFullyQualifiedType, "/")
	fulltyp := pbFullyQualifiedType[idx+1:]

	// 替换RequestType/ResponseType的包名
	idx = strings.LastIndex(fulltyp, ".")
	if idx <= 0 {
		panic(fmt.Errorf("invalid type:%s", fulltyp))
	}

	//pkg := fulltyp[0:idx]
	typ := fulltyp[idx+1:]

	//BUG: https://github.com/hitzhangjie/gorpc-cli/issues/184
	//多个pb文件，package directive相同，但是go_package不同，这种情况下，nfd.Pkg2ValidGoPkg中的package directive
	//到go_package的映射关系是错误的，这里需要具体到文件名才可以。
	//但是，jhump/protoreflect解析出来的rpc的RequestType就是${package-directive}.RequestType这种形式的，没有指明
	//RequestType定义在哪个pb文件中的相关信息。
	//
	//得通过下面的办法，来定位message定义的具体pb文件。
	//
	//- 遍历所有的message及其定义，找到message名称完全匹配的pb文件，找到该pb文件对应的package
	//- 重写请求体的完整类型

	//rtype := pbFullyQualifiedType
	//if gopkg, ok := nfd.Pkg2ValidGoPkg[pkg]; ok && len(gopkg) != 0 {
	//	rtype = gopkg + "." + typ
	//}

	rtype := pbFullyQualifiedType
	if len(nfd.RpcMessageType) != 0 {
		// 找到了message定义的pb文件
		pb, ok := nfd.RpcMessageType[pbFullyQualifiedType]
		if !ok || len(pb) == 0 {
			panic(fmt.Errorf("convert %s to qualified type fail", pbFullyQualifiedType))
		}
		validGoPkg, ok := nfd.Pb2ValidGoPkg[pb]
		if !ok {
			panic(fmt.Errorf("get valid gopkg of %s fail", pb))
		}
		rtype = validGoPkg + "." + typ
	}

	return rtype
}

// PBValidGoPackage return valid go package
func PBValidGoPackage(pkgName string) string {
	var (
		pkg string
	)
	idx := strings.LastIndex(pkgName, "/")
	if idx < 0 {
		pkg = pkgName
	} else {
		pkg = pkgName[idx+1:]
	}

	return strings.Replace(pkg, ".", "_", -1)
}

// Last returns the last element in `list`
func Last(list []string) string {
	idx := len(list) - 1
	return list[idx]
}

func HasPrefix(prefix, str string) bool {
	return strings.HasPrefix(str, prefix)
}

func HasSuffix(prefix, str string) bool {
	return strings.HasSuffix(str, prefix)
}

func Sub(num1, num2 int) int {
	return num1 + num2
}

func LoadGoMod() (mod string, err error) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := filepath.Join(d, "go.mod")
	_, err = os.Lstat(p)
	if err != nil {
		return
	}
	fin, err := os.Open(p)
	if err != nil {
		return
	}
	sc := bufio.NewScanner(fin)
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "module ") {
			return strings.Split(l, " ")[1], nil
		}
	}
	return
}
