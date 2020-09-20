package tpl

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/hitzhangjie/gorpc-cli/util/lang"
)

var funcMap = template.FuncMap{
	"simplify":       lang.PBSimplifyGoType,
	"gopkg":          lang.PBGoPackage,
	"gopkg_simple":   lang.PBValidGoPackage,
	"gotype":         lang.PBGoType,
	"export":         lang.GoExport,
	"gofulltype":     lang.GoFullyQualifiedType,
	"title":          lang.Title,
	"untitle":        lang.UnTitle,
	"trimright":      lang.TrimRight,
	"trimleft":       lang.TrimLeft,
	"splitList":      lang.SplitList,
	"last":           lang.Last,
	"hasprefix":      lang.HasPrefix,
	"hassuffix":      lang.HasSuffix,
	"contains":       strings.Contains,
	"sub":            lang.Sub,
	"camelcase":      strcase.ToCamel,
	"lowercamelcase": strcase.ToLowerCamel,
	"lower":          strings.ToLower,
	"snakecase":      strcase.ToSnake,
}
