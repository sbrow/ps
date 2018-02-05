// +build windows

// Package ps lets you manipulate Adobe Photoshop (CS5) from go.
// This is primarily done by calling VBS/Applescript files.
//
// Currently only works on windows
package ps

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	Cmd  = "cscript.exe"
	Opts = "/nologo"
)

var pkgpath string

func init() {
	_, file, _, _ := runtime.Caller(0)
	pkgpath = path.Dir(file)
}

// Open photoshop.
func Start() error {
	_, err := run("start")
	return err
}

// Open a file.
func Open(path string) ([]byte, error) {
	return run("open", path)
}

// Close the active document.
func Close() error {
	_, err := run("close")
	return err
}

// Quit photoshop with save status.
func Quit(save int) ([]byte, error) {
	return run("quit", string(save))
}

func Js(args ...string) ([]byte, error) {
	return run("dojs", args...)
}
func Wait(msg string) {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
}

func run(name string, args ...string) ([]byte, error) {
	var ext string
	var out bytes.Buffer
	var stderr bytes.Buffer

	switch runtime.GOOS {
	case "windows":
		ext = ".vbs"
	case "darwin":
		ext = ".applescript"
	}
	if !strings.HasSuffix(name, ext) {
		name += ext
	}
	args = append([]string{Opts, path.Join(pkgpath, "scripts", name)}, args...)
	cmd := exec.Command(Cmd, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.Bytes(), err
}
