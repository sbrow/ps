// TODO: Update package tests.
package ps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/sbrow/ps/runner"
)

var testDoc string

func init() {
	Mode = Normal
	log.Printf("Running in mode %v\n", Mode)
	testDoc = filepath.Join(pkgpath, "test.psd")
}

func TestPkgPath(t *testing.T) {
	want := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "sbrow", "ps")
	got := filepath.Join(pkgpath)
	if got != want {
		t.Errorf("wanted: %s\ngot: %s", want, got)
	}
}

func TestInit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestStart\"")
	}
	Quit(2)
	if err := Init(); err != nil {
		t.Error(err)
	}
}

func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestOpen\"")
	}
	if err := Open(testDoc); err != nil {
		log.Println(testDoc)
		t.Fatal(err)
	}
}

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestClose\"")
	}
	err := Close(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestQuit\"")
	}
	Init()
	err := Quit(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoJs(t *testing.T) {
	var err error
	if err = Open(testDoc); err != nil {
		t.Fatal(err)
	}
	want := "F:\\\\TEMP\\\\[0-9]*\r\narg\r\nargs\r\n"
	script := "test.jsx"
	var got []byte
	if got, err = DoJS(script, "arg", "args"); err != nil {
		t.Fatal(err)
	}
	if b, err := regexp.Match(want, got); err != nil || !b {
		fail := fmt.Sprintf("wanted: %s\ngot: %s err: %s\n", want, got, err)
		t.Error(fail)
	}
}

func TestRun(t *testing.T) {
	out := []byte("hello,\r\nworld!\r\n")
	msg, err := runner.Run("test", "hello,", "world!")
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != string(out) {
		fail := fmt.Sprintf("TestRun faild.\ngot:\n\"%s\"\nwant:\n\"%s\"\n", msg, out)
		t.Fatal(fail)
	}
}

func TestWait(t *testing.T) {
	Wait("Waiting...")
}

func TestDoAction_Crop(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestDoAction_Crop\"")
	}
	var err error
	if err = Open(testDoc); err != nil {
		t.Fatal(err)
	}
	if err = DoAction("DK", "Crop"); err != nil {
		t.Error(err)
	}
}

func TestDoAction_Undo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestDoAction_Undo\"")
	}
	if err := DoAction("DK", "Undo"); err != nil {
		t.Fatal(err)
	}
}

func TestSaveAs(t *testing.T) {
	if err := SaveAs("F:\\TEMP\\test.png"); err != nil {
		t.Fatal(err)
	}
	os.Remove("F:\\TEMP\\test.png")
}

func TestLayerSet(t *testing.T) {
	if _, err := NewLayerSet("Group 1/", nil); err != nil {
		t.Fatal(err)
	}
}
func TestMove(t *testing.T) {
	d, err := ActiveDocument()
	if err != nil {
		t.Fatal(err)
	}
	lyr := d.LayerSet("Group 1").ArtLayer("Layer 1")
	lyr.SetPos(100, 50, "TL")
}

func TestActiveDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestDocument\"")
	}
	Open(testDoc)
	d, err := ActiveDocument()
	defer d.Dump()
	if err != nil {
		t.Fatal(err)
	}
	if d != d.artLayers[0].Parent() {
		fmt.Println(d)
		fmt.Println(d.artLayers[0].Parent())
		t.Fatal("ArtLayers do not have doc as parent.")
	}
	if d != d.layerSets[0].Parent() {
		fmt.Println(d)
		fmt.Println(d.layerSets[0].Parent())
		t.Fatal("LayerSets do not have doc as parent.")
	}
	if d.layerSets[0] != d.layerSets[0].artLayers[0].Parent() {
		fmt.Println(d.layerSets[0])
		fmt.Println(d.layerSets[0].artLayers[0])
		fmt.Println(d.layerSets[0].artLayers[0].Parent())
		t.Fatal("Layerset's ArtLayers do not have correct parents")
	}
	lyr := d.LayerSet("Group 1").ArtLayer("Layer 1")
	if lyr == nil {
		t.Fatal("lyr does not exist")
	}
	s := Stroke{Size: 4, Color: &RGB{0, 0, 0}}
	lyr.SetStroke(s, &RGB{128, 128, 128})
}

func TestColor(t *testing.T) {
	byt, err := runner.Run("colorLayer.vbs", "255", "255", "255")
	fmt.Println(string(byt))
	fmt.Println(err)
	if err != nil {

		t.Fatal()
	}
}

func TestApplyDataset(t *testing.T) {
	err := ApplyDataset("Anger")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDocumentLayerSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestDocumentLayerSet")
	}
	d, err := ActiveDocument()
	if err != nil {
		t.Fatal(err)
	}
	set := d.LayerSet("Group 1")
	fmt.Println(set)
	for _, lyr := range set.ArtLayers() {
		fmt.Println(lyr.name)
	}
	lyr := set.ArtLayer("Layer 1")
	fmt.Println(lyr)
	// set = d.LayerSet("Indicators").LayerSet("Life")
	// fmt.Println(set)
	for _, lyr := range set.ArtLayers() {
		fmt.Println(lyr.name)
	}
}

func TestLoadedDoc(t *testing.T) {
	var d *Document
	byt, err := ioutil.ReadFile("./data/Test.psd.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(byt, &d)
	if err != nil {
		t.Fatal(err)
	}
	if d != d.ArtLayers()[0].Parent() {
		t.Fatal("Loaded document's ArtLayers do not point to doc")
	}
	if d != d.LayerSets()[0].Parent() {
		t.Fatal("Loaded document's LayerSets do not point to doc")
	}
	if d.LayerSets()[0] != d.layerSets[0].artLayers[0].Parent() {
		t.Fatal("Loaded document's LayerSet's ArtLayers do not point to layerSets")
	}
}

func TestJSLayer(t *testing.T) {
	d, _ := ActiveDocument()
	set := d.LayerSet("Group 1")
	lyr := set.ArtLayer("Layer 1")
	fmt.Println(JSLayer(set.Path()))
	fmt.Println(JSLayer(lyr.Path()))
}
func TestDoJs_HideLayer(t *testing.T) {
	err := Open(testDoc)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ActiveDocument()
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := NewLayerSet("Group 1/", d)
	lyr.SetVisible(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTextItem(t *testing.T) {
	// err := Open(testDoc)
	// if err != nil {
	// t.Fatal(err)
	// }

	d, err := ActiveDocument()
	if err != nil {
		t.Fatal(err)
	}
	for _, lyr := range d.ArtLayers() {
		if lyr.Name() == "Text" {
			lyr.SetText("Butts")
			// lyr.FmtText(0, 5, "Arial", "Regular")
			// lyr.FmtText(0, 3, "Arial", "Bold")
		}
	}

	/*	byt := []byte(`{"Name": "lyr", "TextItem": {"Contents": "lyr", "Size": 12.000, "Font": "ArialItalic"}}`)
		lyr := &ArtLayer{}
		// byt := []byte(`{"Name": "lyr"}`)
		// lyr := &TextItem{}
		err := lyr.UnmarshalJSON(byt)
		fmt.Printf("%+v\n", lyr)
		fmt.Println(lyr.TextItem)
		if err != nil {
			t.Fatal(err)
		}
	*/
}

func BenchmarkDoc_Go(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ActiveDocument()
		if err != nil {
			b.Fatal(err)
		}
	}
}

//.8s
//.15
func BenchmarkHideLayer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// _, err := Layers("Areas/TitleBackground/")
		// if err != nil {
		// b.Fatal(err)
		// }
	}
}

// 59ns
func BenchmarkHelloWorld_go(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("Hello, world!")
	}
}

// ~35200000ns (.0352s)
func BenchmarkHelloWorld_vbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := runner.Run("helloworld")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ~51700000 (0.0517)
func BenchmarkHelloWorld_js(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := DoJS("test.jsx", "Hello, World!")
		if err != nil {
			b.Fatal(err)
		}
	}
}
