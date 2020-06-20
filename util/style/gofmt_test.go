package style_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hitzhangjie/gorpc/util/style"
)

var f1 string = `package main

func main() {
fmt.Println("hello world")
	fmt.Println("Hello world")
fmt.Println("hello world")
return
}
`

var f2 string = `package main

func main() {
	fmt.Println("hello world")
	fmt.Println("Hello world")
	fmt.Println("hello world")
	return
}
`

func TestGoFmt(t *testing.T) {

	d := filepath.Join(os.TempDir(), "gofmt")
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	// 写入go文件
	f := filepath.Join(d, "f1.go")
	fin, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fin.WriteString(f1)
	if err != nil {
		t.Fatal(err)
	}
	fin.Close()

	// 格式化go文件
	if err := style.GoFmt(f); err != nil {
		t.Fatal(err)
	}

	// 校验go文件
	fin, err = os.Open(f)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ioutil.ReadAll(fin)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != f2 {
		t.Errorf("\n**input**:\n%s\n**output**:\n%s\n**want**:\n%s\n", f1, string(buf), f2)
	}
}

func TestGoFmtDir(t *testing.T) {

	d := filepath.Join(os.TempDir(), "gofmt")
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	// 写入go文件
	f := filepath.Join(d, "f1.go")
	fin, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fin.WriteString(f1)
	if err != nil {
		t.Fatal(err)
	}
	fin.Close()

	// 格式化go文件
	if err := style.GoFmtDir(d); err != nil {
		t.Fatal(err)
	}

	// 校验go文件
	fin, err = os.Open(f)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ioutil.ReadAll(fin)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf) != f2 {
		t.Errorf("\n**input**:\n%s\n**output**:\n%s\n**want**:\n%s\n", f1, string(buf), f2)
	}
}
