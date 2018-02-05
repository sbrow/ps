package ps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	Cmd  = "cscript.exe"
	Opts = "/nologo"
)

var PKGPATH = path.Join(os.Getenv("GOPATH"), "src", "github.com", "sbrow", "ps")

func Start() error {
	_, err := run("start")
	return err
}

func Open(path string) ([]byte, error) {
	return run("open", path)
}

func Close() error {
	_, err := run("close")
	return err
}

func Quit() ([]byte, error) {
	return run("quit")
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
	var dir string
	var out bytes.Buffer
	var stderr bytes.Buffer

	switch runtime.GOOS {
	case "windows":
		ext = ".vbs"
		dir = "win"
	case "darwin":
		ext = ".applescript"
		dir = "mac"
	}
	if !strings.HasSuffix(name, ext) {
		name += ext
	}
	args = append([]string{Opts, path.Join(PKGPATH, dir, name)}, args...)
	cmd := exec.Command(Cmd, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.Bytes(), err
}
