// Package ps is a rudimentary API between Adobe Photoshop and go.
//
// Most of the interaction between the two is implemented via
// javascript and/or VBS/Applescript.
//
// Currently only works with CS5 on Windows.
package ps

import (
	"bytes"
	"encoding/json"
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

// PSSaveOptions is an enum for options when closing a document.
type PSSaveOptions int

func (p *PSSaveOptions) String() string {
	return fmt.Sprint("", *p)
}

const (
	PSSaveChanges         PSSaveOptions = 1
	PSDoNotSaveChanges    PSSaveOptions = 2
	PSPromptToSaveChanges PSSaveOptions = 3
)

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

// Start opens Photoshop.
func Start() error {
	_, err := run("start")
	return err
}

// Close closes the active document.
func Close(save PSSaveOptions) error {
	_, err := run("close", save.String())
	return err
}

// Open opens a file with the specified path.
func Open(path string) error {
	_, err := run("open", path)
	return err
}

// Quit exits Photoshop.
func Quit(save PSSaveOptions) error {
	_, err := run("quit", save.String())
	return err
}

// DoJs runs a Photoshop javascript script file (.jsx) from the specified location.
// It can't directly return output, so instead the scripts write their output to
// a temporary file.
func DoJs(path string, args ...string) ([]byte, error) {
	// Temp file for js to output to.
	outpath := filepath.Join(os.Getenv("TEMP"), "js_out.txt")
	defer os.Remove(outpath)

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
		return cmd, err
	}
	cmd = append(cmd, file...)
	return cmd, err
}

// Wait prints a message to the console and halts operation until the user
// signals that they are ready (by pushing enter).
//
// Useful for when you need to do something by hand in the middle of an
// otherwise automated process.
func Wait(msg string) {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
	fmt.Println()
}

// run handles running the script files, returning output, and displaying errors.
func run(name string, args ...string) ([]byte, error) {
	var ext string
	var out bytes.Buffer
	var errs bytes.Buffer

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
	cmd.Stderr = &errs
	err := cmd.Run()
	if err != nil {
		return out.Bytes(), err
		// return append(out.Bytes(), errs.Bytes()...), err
	}
	if len(errs.Bytes()) != 0 {
		return out.Bytes(), errors.New(string(errs.Bytes()))
	}
	return out.Bytes(), nil
}

// DoAction runs a Photoshop action with name from set.
func DoAction(set, name string) error {
	_, err := run("action", set, name)
	return err
}

// SaveAs saves the Photoshop document file to the given location.
func SaveAs(path string) error {
	_, err := run("save", path)
	return err
}

// Layers returns an array of ArtLayers from the active document
// based on the given path string.
func Layers(path string) ([]ArtLayer, error) {
	byt, err := DoJs("getLayers.jsx", path)
	var out []ArtLayer
	err = json.Unmarshal(byt, &out)
	if err != nil {
		return []ArtLayer{}, err
	}
	return out, err
}

// Layer returns an ArtLayer from the active document given a specified
// path string. Layer calls Layers() and returns the first result.
func Layer(path string) (ArtLayer, error) {
	lyrs, err := Layers(path)
	return lyrs[0], err
}

// ApplyDataset fills out a template file with information from a given dataset (csv) file.
func ApplyDataset(name string) ([]byte, error) {
	return DoJs("applyDataset.jsx", name)
}
