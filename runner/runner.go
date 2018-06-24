//go:generate godocdown -template ../.godocdown.template -output=README.md

// Package runner runs the non-go code that Photoshop understands,
// and passes it to back to the go program. Currently, this is
// primarily implemented through Adobe Extendscript, but hopefully
// in the future it will be upgraded to a C++ plugin.
package runner

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Windows is the runner for Visual Basic Scripts.
var Windows = Runner{
	Cmd:  "cscript.exe",
	Args: []string{"/nologo"},
	Ext:  ".vbs",
}

// pkgpath is the path to this package.
var pkgpath string

// std is the Standard Runner.
var std Runner

func init() {
	_, file, _, _ := runtime.Caller(0)
	pkgpath = filepath.Dir(file)
	switch runtime.GOOS {
	case "windows":
		std = Windows
	}
}

// Runner runs script files to communicate between the OS/Photoshop and the Go code.
type Runner struct {
	Cmd  string   // The name of the command to run
	Args []string // The arguments to pass to the command.
	Ext  string   // The file extension to use for these commands.
}

// Run runs the standard runner with the given values.
func Run(name string, args ...string) ([]byte, error) {
	var out, errs bytes.Buffer
	cmd := exec.Command(std.Cmd, parseArgs(name, args...)...)
	cmd.Stdout, cmd.Stderr = &out, &errs
	if err := cmd.Run(); err != nil || len(errs.Bytes()) != 0 {
		return out.Bytes(), errors.New(errs.String())
	}
	return out.Bytes(), nil
}

// parseArgs parses the given args into the correct syntax.
func parseArgs(name string, args ...string) []string {
	if !strings.HasSuffix(name, std.Ext) {
		name += std.Ext
	}
	newArgs := append(std.Args, filepath.Join(pkgpath, "scripts", name))
	if strings.Contains(name, "dojs") {
		newArgs = append(newArgs, args[0], fmt.Sprint(strings.Join(args[1:], ",,")))
	} else {
		newArgs = append(newArgs, args...)
	}
	return newArgs
}