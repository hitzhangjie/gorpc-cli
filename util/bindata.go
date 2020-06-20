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
var outputdir = flag.String("outputdir", "bindata", "write $file.go under directory $outputdir")

func main() {

	flag.Parse()

	if len(*file) == 0 || len(*outputdir) == 0 {
		panic("invalid arguments: invalid file or outputdir")
	}

	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(fmt.Errorf("read file error: err", err))
	}

	baseName := filepath.Base(*file)
	targetVar := fmt.Sprintf("%sBytes", strcase.ToCamel(baseName))

	tpl := fmt.Sprintf(`package bindata

var %s = []uint8{
{{ range $idx, $val := .Data }}{{$val}},{{ end }}
}
`, targetVar)

	targetBaseName := fmt.Sprintf("static_%s.go", baseName)
	tplInstance, err := template.New(targetBaseName).Parse(tpl)
	if err != nil {
		panic(fmt.Errorf("parse template error: %v", err))
	}

	err = os.MkdirAll(*outputdir, 0777)
	if err != nil {
		panic(fmt.Errorf("create outputdir error: %v", err))
	}

	targetFilePath := filepath.Join(*outputdir, targetBaseName)
	fout, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
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

	fmt.Printf("ok, filedata stored to %s\n", targetBaseName)
}


