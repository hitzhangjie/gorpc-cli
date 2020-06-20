package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
)

var file = flag.String("file", "", "read `file` and store its content into $file.go")

func main() {

	flag.Parse()

	if len(*file) == 0 {
		panic("invalid filename")
	}

	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(fmt.Errorf("read file error: err", err))
	}

	baseName := filepath.Base(*file)
	target := fmt.Sprintf("static_%s.go", baseName)
	targetVar := fmt.Sprintf("%sBytes", strcase.ToCamel(baseName))

	tpl := fmt.Sprintf(`
package bindata

var %s = []uint8{
{{ range $idx, $val := .Data }}{{$val}},{{ end }}
}
`, targetVar)

	tplInstance, err := template.New(target).Parse(tpl)
	if err != nil {
		panic(fmt.Errorf("parse template error: %v", err))
	}

	fout, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Errorf("open file error: %v", err))
	}
	err = tplInstance.Execute(fout, &struct {
		Data []uint8
	}{
		Data: buf,
	})
	if err != nil {
		panic(fmt.Errorf("template execute error: %v", err))
	}

	fmt.Printf("ok, filedata stored to %s\n", target)
}
