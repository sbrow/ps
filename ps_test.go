package ps

import (
	"fmt"
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

/*
func TestOpen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping \"TestOpen\"")
	}
	err := Open("F:\\GitLab\\dreamkeepers-psd\\Template009.1.psd")
	if err != nil {
		t.Fatal(err)
	}
}
*/
/*
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
}

func TestLayers(t *testing.T) {
	l, err := Layers("Areas/TitleBackground/")
	fmt.Println(l)
	if err != nil {
		t.Fatal(err)
	}
}
*/
/*
func TestLayer(t *testing.T) {
	_, err := Layer("Border/Inner Border")
	if err != nil {
		t.Fatal(err)
	}
}*/

/*func TestMove(t *testing.T) {
	lyr, err := Layer("Group 1/Layer 1")
	if err != nil {
		t.Fatal(err)
	}
	lyr.Position(100, 50, "top")
}*/

/*
func TestLayerSet(t *testing.T) {
	set, err := GetLayerSet("Indicators/")
	fmt.Println(set)
	fmt.Println(set.ArtLayers[0].Parent)
	if err != nil {
		t.Fatal(err)
	}
}
*/

func TestDocument(t *testing.T) {
	d, err := GetDocument()
	fmt.Println(d)
	fmt.Println(d.ArtLayers[0])
	fmt.Println(d.ArtLayers[0].Parent)
	fmt.Println(d.LayerSets[0])
	fmt.Println(d.LayerSets[0].Parent)
	fmt.Println(d.LayerSets[0].ArtLayers[0])
	fmt.Println(d.LayerSets[0].ArtLayers[0].Parent)
	fmt.Println(d.LayerSets[0].ArtLayers[0].Parent.Parent())
	if err != nil {
		t.Fatal(err)
	}
	if d != d.ArtLayers[0].Parent {
		t.Fatal("Fucked")
	}
	if d != d.LayerSets[0].Parent() {
		t.Fatal("Fucked")
	}
	if d.LayerSets[0] != d.LayerSets[0].ArtLayers[0].Parent {
		t.Fatal("Fucked")
	}

}

/*func TestActiveDocument(t *testing.T) {
	e, err := DoJs("compilejs.jsx", "alert('testing!')")
	fmt.Println(string(e))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := ActiveDocument()
	fmt.Println(doc)
	if err != nil {
		t.Fatal(err)
	}
}
*/

/*
func TestApplyDataset(t *testing.T) {
	out := []byte("done!\r\n")
	ret, err := ApplyDataset("Anger")
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != string(out) {
		fail := fmt.Sprintf("TestJS failed.\ngot:\t\"%s\"\nwant:\t\"%s\"", ret, out)
		t.Fatal(fail)
	}
	err = Quit(2)
	if err != nil {
		t.Fatal(err)
	}
}
*/

/*func TestDoJs_HideLayer(t *testing.T) {
	_, err := DoJs("setLayerVisibility.jsx", "Areas/TitleBackground", "false")
	if err != nil {
		t.Fatal(err)
	}
}*/

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
