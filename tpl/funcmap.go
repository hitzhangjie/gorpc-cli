package tpl

import (
	"github.com/hitzhangjie/go-rpc-cmdline/parser"
	"text/template"
)

var funcMap = template.FuncMap{
	"simplify":   parser.PBSimplifyGoType,
	"gopkg":      parser.PBGoPackage,
	"gotype":     parser.PBGoType,
	"export":     parser.GoExport,
	"gofulltype": parser.GoFullyQualifiedType,
	"title":      parser.Title,
	"untitle":    parser.UnTitle,
	"trimright":  parser.TrimRight,
	"splitList":  parser.SplitList,
	"last":       parser.Last,
	"hasprefix":  parser.HasPrefix,
}
