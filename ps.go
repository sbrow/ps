// Package ps lets you manipulate Adobe Photoshop (CS5) from go.
// This is primarily done by calling VBS/Applescript files.
//
// Currently only works on windows
package ps

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var Cmd string
var Opts string
var pkgpath string

func init() {
	_, file, _, _ := runtime.Caller(0)
	pkgpath = filepath.Dir(file)
	switch runtime.GOOS {
	case "windows":
		Cmd = "cscript.exe"
		Opts = "/nologo"
	case "darwin":
		Cmd = "osacript"
	}
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

func Js(path string, args ...string) ([]byte, error) {
	// Temp file for js to output to.
	outpath := filepath.Join(os.Getenv("TEMP"), "js_out.txt")
	args = append([]string{outpath}, args...)

	// If passed a script by name, assume it's in the default folder.
	if filepath.Dir(path) == "." {
		path = filepath.Join(pkgpath, "scripts", path)
	}

	args = append([]string{path}, args...)
	cmd, err := run("dojs", args...)
	if err != nil {
		return []byte{}, err
	}
	file, err := ioutil.ReadFile(outpath)
	if err != nil {
		// fmt.Println(cmd)
		return cmd, err
	}
	cmd = append(cmd, file...)
	// os.Remove(outpath)
	return cmd, err
}

func Wait(msg string) {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
}

func run(name string, args ...string) ([]byte, error) {
	var ext string
	var out bytes.Buffer

	switch runtime.GOOS {
	case "windows":
		ext = ".vbs"
	case "darwin":
		ext = ".applescript"
	}
	if !strings.HasSuffix(name, ext) {
		name += ext
	}

	if strings.Contains(name, "dojs") {
		args = append([]string{Opts, filepath.Join(pkgpath, "scripts", name)},
			args[0],
			fmt.Sprintf("%s", strings.Join(args[1:], ",")),
		)
	} else {
		args = append([]string{Opts, filepath.Join(pkgpath, "scripts", name)}, args...)
	}

	cmd := exec.Command(Cmd, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return []byte{}, errors.New(string(out.Bytes()))
	} else {
		return out.Bytes(), err
	}
}
