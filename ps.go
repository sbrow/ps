// Package ps is a rudimentary API between Adobe Photoshop CS5 and Golang.
// The interaction between the two is implemented using Javascript/VBScript.
//
// Currently only supports Photoshop CS5 Windows x86_64.
//
// TODO: Creatue a Photoshop struct to hold program values and functions.
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
	// "update"
)

// The name of the program that runs scripts on this OS.
var scmd string

// The options we need to pass to scmd.
var opts string
var pkgpath string

func init() {
	_, file, _, _ := runtime.Caller(0)
	pkgpath = filepath.Dir(file)
	switch runtime.GOOS {
	case "windows":
		scmd = "cscript.exe"
		opts = "/nologo"
	case "darwin":
		scmd = "osacript"
	}
	// update.Update()
}

// ApplyDataset fills out a template file with information
// from a given dataset (csv) file. It is important to note that running this
// function will change data in the Photoshop document, but will not update
// data in the Go Document struct- you will have to implement syncing
// them yourself.
func ApplyDataset(name string) error {
	_, err := DoJs("applyDataset.jsx", name)
	return err
}

// Close closes the active document in Photoshop, using the given save option.
func Close(save PSSaveOptions) error {
	_, err := run("close", save.String())
	return err
}

// DoAction runs the Photoshop Action "name" from the Action Set "set".
func DoAction(set, name string) error {
	_, err := run("action", set, name)
	return err
}

// DoJs runs a Photoshop Javascript script file (.jsx) from the specified location.
// The script can't directly return output, so instead it writes output to
// a temporary file ($TEMP/js_out.txt), whose contents is then read and returned.
func DoJs(path string, args ...string) (out []byte, err error) {
	// Temp file for js to output to.
	outpath := filepath.Join(os.Getenv("TEMP"), "js_out.txt")
	// defer os.Remove(outpath)
	if !strings.HasSuffix(path, ".jsx") {
		path += ".jsx"
	}

	args = append([]string{outpath}, args...)

	// If passed a script by name, assume it's in the default folder.
	if filepath.Dir(path) == "." {
		path = filepath.Join(pkgpath, "scripts", path)
	}

	args = append([]string{path}, args...)
	cmd, err := run("dojs", args...)
	if err == nil {
		file, err := ioutil.ReadFile(outpath)
		if err == nil {
			cmd = append(cmd, file...)
		}
	}
	return cmd, err
}

// JSLayer "compiles" Javascript code to get an ArtLayer with the given path.
// The output always ends with a semicolon, so if you want to access a specific
// property of the layer, you'll have to trim the output before concatenating.
//
// TODO: get rid of the semicolon.
func JSLayer(path string, art ...bool) string {
	path = strings.TrimLeft(path, "/")
	pth := strings.Split(path, "/")
	js := "app.activeDocument"
	last := len(pth) - 1
	if len(art) > 0 {
		pth = pth[:len(pth)-1]
		last--
	}
	if last > 0 {
		for i := 0; i < last; i++ {
			js += fmt.Sprintf(".layerSets.getByName('%s')", pth[i])
		}
	}
	if pth[last] != "" {
		js += fmt.Sprintf(".artLayers.getByName('%s')", pth[last])
	}
	return js + ";"
}

// Open opens a Photoshop document with the specified path.
// If Photoshop is not currently running, it is started before
// opening the document.
func Open(path string) error {
	_, err := run("open", path)
	return err
}

// Quit exits Photoshop using the given save option.
func Quit(save PSSaveOptions) error {
	_, err := run("quit", save.String())
	return err
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
		args = append([]string{opts, filepath.Join(pkgpath, "scripts", name)},
			args[0],
			fmt.Sprintf("%s", strings.Join(args[1:], ",,")),
		)
	} else {
		args = append([]string{opts, filepath.Join(pkgpath, "scripts", name)}, args...)
	}
	cmd := exec.Command(scmd, args...)
	cmd.Stdout = &out
	cmd.Stderr = &errs
	err := cmd.Run()
	if err != nil || len(errs.Bytes()) != 0 {
		return out.Bytes(), errors.New(string(errs.Bytes()))
	}
	return out.Bytes(), nil
}

// SaveAs saves the Photoshop document to the given location.
func SaveAs(path string) error {
	_, err := run("save", path)
	return err
}

// Start opens Photoshop.
func Start() error {
	_, err := run("start")
	return err
}

// Wait prints a message to the console and halts operation until the user
// signals that they are ready to continue (by pushing enter).
//
// Useful for when you need to do something by hand in the middle of an
// otherwise automated process. (i.e. importing a dataset).
func Wait(msg string) {
	fmt.Print(msg)
	var input string
	fmt.Scanln(&input)
	fmt.Println()
}
