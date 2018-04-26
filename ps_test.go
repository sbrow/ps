package ps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPkgPath(t *testing.T) {
	out := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "sbrow", "ps")
	if filepath.Join(pkgpath) != out {
		t.Fatal(filepath.Join(pkgpath), out)
	}
}

func TestStart(t *testing.T) {
	err := Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpen(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("Skipping \"TestOpen\"")
	// }
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
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
	err := Quit(2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoJs(t *testing.T) {
	out := []byte("F:\\TEMP\\js_out.txt\r\narg\r\nargs\r\n")
	script := "test.jsx"
	ret, err := DoJs(script, "arg", "args")
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != string(out) {
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", ret, out)
		t.Fatal(fail)
	}
}

func TestRun(t *testing.T) {
	out := []byte("hello,\r\nworld!\r\n")
	msg, err := run("test", "hello,", "world!")
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
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
	err = DoAction("DK", "Crop")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoAction_Undo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestDoAction_Undo\"")
	}
	err := DoAction("DK", "Undo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveAs(t *testing.T) {
	err := SaveAs("F:\\TEMP\\test.png")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("F:\\TEMP\\test.png")
}

func TestLayerSet(t *testing.T) {
	_, err := NewLayerSet("Areas/TitleBackground/", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLayer(t *testing.T) {
	_, err := layer("Border/Inner Border")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMove(t *testing.T) {
	lyr, err := layer("Group 1/Layer 1")
	if err != nil {
		t.Fatal(err)
	}
	lyr.SetPos(100, 50, "TL")
}

func TestActiveDocument(t *testing.T) {
	Mode = Safe
	if testing.Short() {
		t.Skip("Skipping \"TestDocument\"")
	}
	d, err := ActiveDocument()
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
	// d.LayerSet("Areas").LayerSet("Bottom").ArtLayer("L Bar").SetColor(155, 255, 255)
	lyr := d.LayerSet("Text").ArtLayer("speed")
	if lyr == nil {
		t.Fatal("lyr does not exist")
	}
	s := Stroke{Size: 4, Color: &RGB{0, 0, 0}}
	lyr.SetStroke(s, &RGB{128, 128, 128})
	d.Dump()
}

func TestColor(t *testing.T) {
	byt, err := run("colorLayer.vbs", "255", "255", "255")
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
	set := d.LayerSet("Text")
	fmt.Println(set)
	for _, lyr := range set.ArtLayers() {
		fmt.Println(lyr.name)
	}
	lyr := set.ArtLayer("id")
	fmt.Println(lyr)
	set = d.LayerSet("Indicators").LayerSet("Life")
	fmt.Println(set)
	for _, lyr := range set.ArtLayers() {
		fmt.Println(lyr.name)
	}
}

func TestLoadedDoc(t *testing.T) {
	var d *Document
	byt, err := ioutil.ReadFile("Document.txt")
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

func TestDoJs_HideLayer(t *testing.T) {
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
	lyr, err := NewLayerSet("Areas/TitleBackground", nil)
	lyr.SetVisible(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTextItem(t *testing.T) {
	// err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
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
			lyr.FmtText(0, 5, "Arial", "Regular")
			lyr.FmtText(0, 3, "Arial", "Bold")
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
		fmt.Sprintf("Hello, world!")
	}
}

// ~35200000ns (.0352s)
func BenchmarkHelloWorld_vbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := run("helloworld")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ~51700000 (0.0517)
func BenchmarkHelloWorld_js(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := DoJs("test.jsx", "Hello, World!")
		if err != nil {
			b.Fatal(err)
		}
	}
}
