//go:generate sh -c "godoc2md -template ./.doc.template github.com/sbrow/ps > README.md"

// Package ps is a rudimentary API between Adobe Photoshop CS5 and Golang.
// The interaction between the two is implemented using Javascript/VBScript.
//
// Use it to control Photoshop, edit documents, and perform batch operations.
//
// Currently only supports Photoshop CS5 Windows x86_64.
package ps

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/sbrow/ps/runner"
)

// The full path to this directory.
var pkgpath string

func init() {
	_, file, _, _ := runtime.Caller(0)
	pkgpath = filepath.Dir(file)
}

// ApplyDataset fills out a template file with information
// from a given dataset (csv) file. It's important to note that running this
// function will change data in the Photoshop document, but will not update
// data in the Go Document struct- you will have to implement syncing
// them yourself.
func ApplyDataset(name string) error {
	_, err := DoJS("applyDataset.jsx", name)
	return err
}

// Close closes the active document in Photoshop, using the given save option.
// TODO(sbrow): refactor Close to Document.Close
func Close(save SaveOption) error {
	_, err := runner.Run("close", fmt.Sprint(save))
	return err
}

// DoAction runs the Photoshop Action with the given name from the Action Set "from".
func DoAction(action, from string) error {
	_, err := runner.Run("action", action, from)
	return err
}

// DoJS runs a Photoshop Javascript script file (.jsx) from the specified location.
// The script can't directly return output, so instead it writes output to
// a temporary file ($TEMP/js_out.txt), whose contents is then read and returned.
func DoJS(path string, args ...string) ([]byte, error) {
	// var err error
	// Temp file for js to output to.
	outpath, err := ioutil.TempFile(os.Getenv("TEMP"), "")
	defer func() {
		if err = outpath.Close(); err != nil {
			log.Println(err)
		}
		if err = os.Remove(outpath.Name()); err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(path, ".jsx") {
		path += ".jsx"
	}

	args = append([]string{outpath.Name()}, args...)

	// If passed a script by name, assume it's in the default folder.
	if filepath.Dir(path) == "." {
		path = filepath.Join(pkgpath, "runner", "scripts", path)
	}

	args = append([]string{path}, args...)
	cmd, err := runner.Run("dojs", args...)
	if err == nil {
		var data []byte
		data, err = ioutil.ReadFile(outpath.Name())
		if err == nil {
			cmd = append(cmd, data...)
		}
	}
	return cmd, err
}

// Init opens Photoshop if it is not open already.
//
// Init should be called before all other
func Init() error {
	_, err := runner.Run("start")
	return err
}

// JSLayer "compiles" Javascript code to get an ArtLayer with the given path.
// The output always ends with a semicolon, so if you want to access a specific
// property of the layer, you'll have to trim the output before concatenating.
//
// TODO(sbrow): get rid of the semicolon at the end of JSLayer.
func JSLayer(path string) string {
	pth := strings.Split(path, "/")
	js := "app.activeDocument"
	last := len(pth) - 1
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

// JSLayerMerge gets the Javascript code to get the Layer or LayerSet with this path
// and returns the result if you were to merge the bottom-most LayerSet.
//
// If the bottom-most Object in the path is not a LayerSet, it will returns the same
// results as JSLayer.
func JSLayerMerge(path string) string {
	reg := regexp.MustCompile(`layerSets(\.getByName\('[^']*'\)($|[^.]))`)
	return reg.ReplaceAllString(JSLayer(path), "artLayers$1")
}

// Open opens a Photoshop document with the specified path.
// If Photoshop is not currently running, it is started before
// opening the document.
func Open(path string) error {
	_, err := runner.Run("open", path)
	return err
}

// Quit exits Photoshop, closing all open documents using the given save option.
func Quit(save SaveOption) error {
	_, err := runner.Run("quit", fmt.Sprint(save))
	return err
}

// SaveAs saves the Photoshop document to the given location.
func SaveAs(path string) error {
	_, err := runner.Run("save", path)
	return err
}

// Wait prints a message to the console and halts operation until the user
// signals that they are ready to continue (by pushing enter).
//
// Useful for when you need to do something by hand in the middle of an
// otherwise automated process, (e.g. importing a dataset).
func Wait(msg string) {
	var input string

	fmt.Print(msg)
	fmt.Scanln(&input)
	fmt.Println()
}
