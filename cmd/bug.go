/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/hitzhangjie/gorpc-cli/util/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bugCmd)
}

// bugCmd represents the bug command
var bugCmd = &cobra.Command{
	Use:   "bug",
	Short: config.LoadTranslation("bugCmdUsage", nil),
	Long:  config.LoadTranslation("bugCmdUsageLong", nil),
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO make params `project`, `body` configurable
		project := "hitzhangjie/gorpc-cli"
		body := generateBugBody()

		url := fmt.Sprintf("https://github.com/%s/issues/new?body=%s", project, url.QueryEscape(body))
		if !browser.Open(url) {
			fmt.Printf("Please file a new issue at github.com/%s/issues/new using this template:\n\n", project)
			fmt.Print(body)
		}
		return nil
	},
}

// a bug/issue body should conform to this template:
func generateBugBody() string {

	buf := bytes.Buffer{}

	// bug header
	buf.WriteString(bugHeader)

	// gorpc version
	printGoRPCVersion(&buf)

	// question of latest release
	buf.WriteString("### Does this issue reproduce with the latest release?\n\n\n")

	// go env version
	printGoEnvDetails(&buf)

	// footer
	buf.WriteString(bugFooter)

	return buf.String()
}

func printGoRPCVersion(w io.Writer) {
	fmt.Fprintf(w, "### What version of gorpc-cli are you using (`gorpc version`)?\n\n")
	fmt.Fprintf(w, "<pre>\n")
	fmt.Fprintf(w, "$ gorpc version\n")
	printCmdOut(w, "", "gorpc", "version")
	fmt.Fprintf(w, "</pre>\n")
	fmt.Fprintf(w, "\n")
}

// printCmdOut prints the output of running the given command.
// It ignores failures; 'go bug' is best effort.
func printCmdOut(w io.Writer, prefix, path string, args ...string) {
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s%s\n", prefix, bytes.TrimSpace(out))
}

func printGoEnvDetails(w io.Writer) {
	fmt.Fprintf(w, "### What Go configurations, OS or Processor are you using (`go env`)?\n\n")
	fmt.Fprintf(w, "<details><summary><code>go env</code> Output</summary><br><pre>\n")
	fmt.Fprintf(w, "$ go env\n")
	printCmdOut(w, "", "go", "env")
	printGoDetails(w)
	printOSDetails(w)
	printCDetails(w)
	fmt.Fprintf(w, "</pre></details>\n\n")
}

func printGoDetails(w io.Writer) {
	printCmdOut(w, "GOROOT/bin/go version: ", filepath.Join(runtime.GOROOT(), "bin/go"), "version")
	printCmdOut(w, "GOROOT/bin/go tool compile -V: ", filepath.Join(runtime.GOROOT(), "bin/go"), "tool", "compile", "-V")
}

func printOSDetails(w io.Writer) {
	switch runtime.GOOS {
	case "darwin", "ios":
		printCmdOut(w, "uname -v: ", "uname", "-v")
		printCmdOut(w, "", "sw_vers")
	case "linux":
		printCmdOut(w, "uname -sr: ", "uname", "-sr")
		printCmdOut(w, "", "lsb_release", "-a")
		printGlibcVersion(w)
	case "openbsd", "netbsd", "freebsd", "dragonfly":
		printCmdOut(w, "uname -v: ", "uname", "-v")
	case "illumos", "solaris":
		// Be sure to use the OS-supplied uname, in "/usr/bin":
		printCmdOut(w, "uname -srv: ", "/usr/bin/uname", "-srv")
		out, err := ioutil.ReadFile("/etc/release")
		if err == nil {
			fmt.Fprintf(w, "/etc/release: %s\n", out)
		} else {
			panic(fmt.Errorf("failed to read /etc/release: %v", err))
		}
	}
}

func printCDetails(w io.Writer) {
	printCmdOut(w, "lldb --version: ", "lldb", "--version")
	cmd := exec.Command("gdb", "--version")
	out, err := cmd.Output()
	if err == nil {
		// There's apparently no combination of command line flags
		// to get gdb to spit out its version without the license and warranty.
		// Print up to the first newline.
		fmt.Fprintf(w, "gdb --version: %s\n", firstLine(out))
	} else {
		panic(fmt.Errorf("failed to run gdb --version: %v", err))
	}
}

// firstLine returns the first line of a given byte slice.
func firstLine(buf []byte) []byte {
	idx := bytes.IndexByte(buf, '\n')
	if idx > 0 {
		buf = buf[:idx]
	}
	return bytes.TrimSpace(buf)
}

// printGlibcVersion prints information about the glibc version.
// It ignores failures.
func printGlibcVersion(w io.Writer) {
	tempdir := os.TempDir()
	if tempdir == "" {
		return
	}
	src := []byte(`int main() {}`)
	srcfile := filepath.Join(tempdir, "go-bug.c")
	outfile := filepath.Join(tempdir, "go-bug")
	err := ioutil.WriteFile(srcfile, src, 0644)
	if err != nil {
		return
	}
	defer os.Remove(srcfile)
	cmd := exec.Command("gcc", "-o", outfile, srcfile)
	if _, err = cmd.CombinedOutput(); err != nil {
		return
	}
	defer os.Remove(outfile)

	cmd = exec.Command("ldd", outfile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	re := regexp.MustCompile(`libc\.so[^ ]* => ([^ ]+)`)
	m := re.FindStringSubmatch(string(out))
	if m == nil {
		return
	}
	cmd = exec.Command(m[1])
	out, err = cmd.Output()
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s: %s\n", m[1], firstLine(out))

	// print another line (the one containing version string) in case of musl libc
	if idx := bytes.IndexByte(out, '\n'); bytes.Index(out, []byte("musl")) != -1 && idx > -1 {
		fmt.Fprintf(w, "%s\n", firstLine(out[idx+1:]))
	}
}

const (
	bugHeader = `<!-- Please answer these questions before submitting your issue. Thanks! -->

`
	bugFooter = `### What did you do?

<!--
If possible, provide a recipe for reproducing the error.
A complete runnable program is good.
A link on play.golang.org is best.
-->



### What did you expect to see?



### What did you see instead?

`
)
